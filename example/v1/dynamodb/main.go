package main

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
)

func main() {
	dash := machinery.NewDashboard(&config.Config{
		ResultBackend: "http://localhost:8000",
		DynamoDB: &config.DynamoDBConfig{
			TaskStatesTable: "task_states",
			GroupMetasTable: "group_metas",
		},
	})

	dash.ViewAllDeadJobs()
}
