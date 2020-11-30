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

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("asia"),
		Endpoint: aws.String("http://localhost:8000")})
	if err != nil {
		log.Fatal(err)
	}
	dyn := dynamodb.New(sess)
	dash := machinery.NewDashboard(&config.Config{
		ResultBackend: "http://localhost:8000",
		DynamoDB: &config.DynamoDBConfig{
			TaskStatesTable: "task_states",
			GroupMetasTable: "group_metas",
			Client:          dyn,
		},
	})

	res, _ := dash.FindAllTasksByState(tasks.StateFailure)
	bt, _ := json.Marshal(res)
	fmt.Println(string(bt))
}
