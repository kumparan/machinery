package tasks

import (
	"fmt"
	"github.com/RichardKnop/machinery/v1/utils"
	"time"

	"github.com/google/uuid"
)

// Arg represents a single argument passed to invocation fo a task
type Arg struct {
	Name  string      `bson:"name"`
	Type  string      `bson:"type"`
	Value interface{} `bson:"value"`
}

// Headers represents the headers which should be used to direct the task
type Headers map[string]interface{}

// Set on Headers implements opentracing.TextMapWriter for trace propagation
func (h Headers) Set(key, val string) {
	h[key] = val
}

// ForeachKey on Headers implements opentracing.TextMapReader for trace propagation.
// It is essentially the same as the opentracing.TextMapReader implementation except
// for the added casting from interface{} to string.
func (h Headers) ForeachKey(handler func(key, val string) error) error {
	for k, v := range h {
		// Skip any non string values
		stringValue, ok := v.(string)
		if !ok {
			continue
		}

		if err := handler(k, stringValue); err != nil {
			return err
		}
	}

	return nil
}

// Signature represents a single task invocation
type Signature struct {
	UUID           string       `json:"UUID,omitempty"`
	Name           string       `json:"Name,omitempty"`
	RoutingKey     string       `json:"RoutingKey,omitempty"`
	ETA            *time.Time   `json:"ETA,omitempty"`
	GroupUUID      string       `json:"GroupUUID,omitempty"`
	GroupTaskCount int          `json:"GroupTaskCount,omitempty"`
	Args           []Arg        `json:"Args,omitempty"`
	Headers        Headers      `json:"Headers,omitempty"`
	Priority       uint8        `json:"Priority,omitempty"`
	Immutable      bool         `json:"Immutable,omitempty"`
	RetryCount     int          `json:"RetryCount,omitempty"`
	RetryTimeout   int          `json:"RetryTimeout,omitempty"`
	OnSuccess      []*Signature `json:"OnSuccess,omitempty"`
	OnError        []*Signature `json:"OnError,omitempty"`
	ChordCallback  *Signature   `json:"ChordCallback,omitempty"`
	//MessageGroupId for Broker, e.g. SQS
	BrokerMessageGroupId string `json:"BrokerMessageGroupId,omitempty"`
	//ReceiptHandle of SQS Message
	SQSReceiptHandle string
	// StopTaskDeletionOnError used with sqs when we want to send failed messages to dlq,
	// and don't want machinery to delete from source queue
	StopTaskDeletionOnError bool
	// IgnoreWhenTaskNotRegistered auto removes the request when there is no handeler available
	// When this is true a task with no handler will be ignored and not placed back in the queue
	IgnoreWhenTaskNotRegistered bool `json:"IgnoreWhenTaskNotRegistered,omitempty"`
}

// NewSignature creates a new task signature
func NewSignature(name string, args []Arg) (*Signature, error) {
	signatureID := uuid.New().String()
	return &Signature{
		UUID: fmt.Sprintf("task_%v", signatureID),
		Name: name,
		Args: args,
	}, nil
}

func CopySignatures(signatures ...*Signature) []*Signature {
	var sigs = make([]*Signature, len(signatures))
	for index, signature := range signatures {
		sigs[index] = CopySignature(signature)
	}
	return sigs
}

func CopySignature(signature *Signature) *Signature {
	var sig = new(Signature)
	_ = utils.DeepCopy(sig, signature)
	return sig
}
