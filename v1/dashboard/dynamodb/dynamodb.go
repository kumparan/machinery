package dashboard

import (
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Dashboard monitor tasks
type Dashboard struct {
	cnf    *config.Config
	client dynamodbiface.DynamoDBAPI
}

// New :nodoc:
func New(cnf *config.Config) *Dashboard {
	dash := &Dashboard{cnf: cnf}
	if cnf.DynamoDB != nil && cnf.DynamoDB.Client != nil {
		dash.client = cnf.DynamoDB.Client
	} else if cnf.ResultBackend != "" {
		sess := session.Must(session.NewSession(
			&aws.Config{
				Region:   aws.String("asia"),
				Endpoint: aws.String("http://localhost:8000"),
			}),
		)
		dash.client = dynamodb.New(sess)
	} else {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		dash.client = dynamodb.New(sess)
	}

	return dash
}

// FindAllTasksByState :nodoc:
func (m *Dashboard) FindAllTasksByState(state string) (taskStates []*tasks.TaskState, err error) {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(m.cnf.DynamoDB.TaskStatesTable),
		IndexName:              aws.String(tasks.TaskStateIndex), // use secondary global index
		KeyConditionExpression: aws.String("#st = :st"),
		ExpressionAttributeNames: map[string]*string{
			"#st":  aws.String("State"),
			"#err": aws.String("Error"),
		},
		ProjectionExpression: aws.String("TaskName, #err, Signature"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":st": {
				S: aws.String(state),
			},
		},
	}

	out, err := m.client.Query(queryInput)
	if err != nil {
		log.ERROR.Print(err)
		return
	}

	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &taskStates)
	if err != nil {
		log.ERROR.Print(err)
		return
	}

	return
}
