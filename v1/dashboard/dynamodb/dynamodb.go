package dashboard

import (
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Dashboard monitor tasks
type Dashboard struct {
	cnf    *config.Config
	client dynamodbiface.DynamoDBAPI
}

// New :nodoc:
func New(cnf *config.Config) *Dashboard {
	dash := &Dashboard{}
	if cnf.DynamoDB != nil && cnf.DynamoDB.Client != nil {
		dash.client = cnf.DynamoDB.Client
	} else {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		dash.client = dynamodb.New(sess)
	}

	return dash
}

// ViewAllDeadJobs :nodoc:
func (m *Dashboard) ViewAllDeadJobs() {
	table := m.cnf.DynamoDB.TaskStatesTable
	input := &dynamodb.ScanInput{
		TableName: aws.String(table),
	}

	res, err := m.client.Scan(input)
	if err != nil {
		return
	}

	for _, item := range res.Items {
		log.DEBUG.Println(item)
	}
}
