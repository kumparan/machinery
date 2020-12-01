package dashboard

import (
	"github.com/RichardKnop/machinery/v1/backends/result"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type server interface {
	SendTask(signature *tasks.Signature) (*result.AsyncResult, error)
}

// Dashboard monitor tasks
type Dashboard struct {
	cnf    *config.Config
	client dynamodbiface.DynamoDBAPI
	server server
}

// TaskWithSignature :nodoc:
type TaskWithSignature struct {
	TaskName  string `bson:"task_name"`
	Signature string `bson:"signature"`
	Error     string `bson:"error"`
}

// New :nodoc:
func New(cnf *config.Config, srv server) *Dashboard {
	dash := &Dashboard{
		cnf:    cnf,
		server: srv,
	}

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

// FindAllTasksByState :nodoc:
func (m *Dashboard) FindAllTasksByState(state string) (taskStates []*TaskWithSignature, err error) {
	var cursor map[string]*dynamodb.AttributeValue
	var items []map[string]*dynamodb.AttributeValue

	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(m.cnf.DynamoDB.TaskStatesTable),
		IndexName:              aws.String(tasks.TaskStateIndex), // use secondary global index
		Limit:                  aws.Int64(10),
		ProjectionExpression:   aws.String("TaskName, #err, Signature"),
		KeyConditionExpression: aws.String("#st = :st"),
		ExpressionAttributeNames: map[string]*string{
			"#st":  aws.String("State"),
			"#err": aws.String("Error"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":st": {
				S: aws.String(state),
			},
		},
	}

	for {
		queryInput.ExclusiveStartKey = cursor
		out, err := m.client.Query(queryInput)
		if err != nil {
			log.ERROR.Print(err)
			return nil, err
		}

		items = append(items, out.Items...)
		cursor = out.LastEvaluatedKey

		if out.LastEvaluatedKey == nil || len(out.Items) == 0 {
			break
		}
	}

	err = dynamodbattribute.UnmarshalListOfMaps(items, &taskStates)
	if err != nil {
		log.ERROR.Print(err)
		return nil, err
	}

	return
}

// ReEnqueueTask :FIXME: failed to enqueue because the args value not matching the type
func (m *Dashboard) ReEnqueueTask(sig *tasks.Signature) error {
	sig.UUID = ""
	sig.ETA = nil

	_, err := m.server.SendTask(sig)
	return err
}
