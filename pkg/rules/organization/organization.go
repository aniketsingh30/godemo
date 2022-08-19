package maturity

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	Validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"riscvue.com/pkg/entities/customer"
)

type Rules struct{}

func NewRules() *Rules {
	return &Rules{}
}

func (r *Rules) Migrate(connection *dynamodb.DynamoDB) error {
	return r.createTable(connection)
}

func (r *Rules) createTable(connection *dynamodb.DynamoDB) error {
	table := &customer.Customer{}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("MaturityFWId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("MaturityFWId"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(table.TableName()),
	}
	response, err := connection.CreateTable(input)
	if err != nil && strings.Contains(err.Error(), "Table already exists") {
		return nil
	}
	if response != nil && strings.Contains(response.GoString(), "TableStatus: \"CREATING\"") {
		time.Sleep(3 * time.Second)
		err = r.createTable(connection)
		if err != nil {
			return err
		}
	}
	return err
}

func (r *Rules) GetMock() interface{} {
	return customer.Customer{

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (r *Rules) Validate(model interface{}) error {

	customeModel, err := customer.InterfaceToModel(model)
	if err != nil {
		return err
	}
	return Validation.ValidateStruct(customeModel,

		Validation.Field(&customeModel.CustomerName, Validation.Required, Validation.Length(3, 50)),
		Validation.Field(&customeModel.CreatedBy, Validation.Required, Validation.Length(3, 50)),
		Validation.Field(&customeModel.OwnerId, Validation.Required, Validation.Length(3, 50)),
	)
}

func (r *Rules) ValidateUpdate(model interface{}) error {

	customeModel, err := customer.InterfaceToModelUpdate(model)
	if err != nil {
		return err
	}
	return Validation.ValidateStruct(customeModel,
		Validation.Field(&customeModel.CustomerId, Validation.Required, is.UUIDv4),
		Validation.Field(&customeModel.CustomerName, Validation.Required, Validation.Length(3, 50)),
		Validation.Field(&customeModel.FieldName, Validation.Required, Validation.Length(3, 50)),

		Validation.Field(&customeModel.UpdatedBy, Validation.Required, Validation.Length(3, 50)),
	)
}

func (r *Rules) ConvertIoReaderToStruct(data io.Reader, model interface{}) (interface{}, error) {
	if data == nil {
		return nil, errors.New("body is invalid")
	}
	return model, json.NewDecoder(data).Decode(model)
}
