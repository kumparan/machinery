package iface

import "github.com/RichardKnop/machinery/v1/tasks"

// Dashboard :noodc:
type Dashboard interface {
	FindAllTasksByState(state string) (taskStates []*tasks.TaskState, err error)
}
