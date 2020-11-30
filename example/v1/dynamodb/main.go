package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var sigj = `{"Name":"handle_dead_letter_queue","RoutingKey":"commerce-service-dlq-worker","ETA":"2020-11-30T08:36:06.921204Z","Args":[{"Name":"","Type":"string","Value":"receipt not found"},{"Name":"task","Type":"string","Value":"create_user_subscription"},{"Name":"payload","Type":"string","Value":"-99"}],"RetryTimeout":8}`

func main() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("asia"),
		Endpoint: aws.String("http://localhost:8000"),
	}))

	cfg := &config.Config{
		Broker:        "redis://localhost:6379/3",
		ResultBackend: "http://localhost:8000",
		DynamoDB: &config.DynamoDBConfig{
			TaskStatesTable: "task_states",
			GroupMetasTable: "group_metas",
			Client:          dynamodb.New(sess),
		},
	}

	fmt.Println("init dashboard")
	dash, err := machinery.NewDashboard(cfg)
	if err != nil {
		log.Fatal(err)
	}

	res, _ := dash.FindAllTasksByState(tasks.StateFailure)
	bt, _ := json.Marshal(res)
	fmt.Println(string(bt))

	fmt.Print("\n\n================= ReEnqueueTask =================\n\n")

	var sig tasks.Signature
	json.Unmarshal([]byte(sigj), &sig)
	dash.ReEnqueueTask(&sig)
}
