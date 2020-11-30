package iface

import (
	dynamodbdash "github.com/RichardKnop/machinery/v1/dashboard/dynamodb"
	"github.com/RichardKnop/machinery/v1/tasks"
)

// Dashboard :noodc:
type Dashboard interface {
	FindAllTasksByState(state string) (taskStates []*dynamodbdash.TaskWithSignature, err error)
	ReEnqueueTask(sig *tasks.Signature) error
}
