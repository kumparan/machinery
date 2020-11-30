package machinery

import (
	"github.com/RichardKnop/machinery/v1/config"
	dashboardiface "github.com/RichardKnop/machinery/v1/dashboard/iface"
)

// Dashboard :nodoc:
type Dashboard struct {
	cnf *config.Config
}

// NewDashboard :nodoc:
func NewDashboard(cnf *config.Config) dashboardiface.Dashboard {
	dash, _ := DashboardFactory(cnf)
	return dash
}
