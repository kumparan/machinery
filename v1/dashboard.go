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
func NewDashboard(cnf *config.Config) (dashboardiface.Dashboard, error) {
	srv, err := NewServer(cnf)
	if err != nil {
		return nil, err
	}

	dash, err := DashboardFactory(cnf, srv)
	return dash, err
}
