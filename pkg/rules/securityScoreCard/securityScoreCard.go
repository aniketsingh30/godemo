package securityScoreCard

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
	"github.com/google/uuid"
	"riscvue.com/pkg/entities"
	"riscvue.com/pkg/entities/securityScoreCard"
)

type Rules struct{}

func NewRules() *Rules {
	return &Rules{}
}

func (r *Rules) Migrate(connection *dynamodb.DynamoDB) error {
	return r.createTable(connection)
}

func (r *Rules) createTable(connection *dynamodb.DynamoDB) error {
	table := &securityScoreCard.MainResponse{}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
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
	return securityScoreCard.MainResponse{
		Base: entities.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func (r *Rules) Validate(model interface{}) error {

	assessmentModel, err := securityScoreCard.InterfaceToModel(model)
	if err != nil {
		return err
	}
	return Validation.ValidateStruct(assessmentModel,
		//	Validation.Field(&assessmentModel.ID, Validation.Required, is.UUIDv4),
		Validation.Field(&assessmentModel.DomainName, Validation.Required, Validation.Length(3, 50)),
		Validation.Field(&assessmentModel.CreatedBy, Validation.Required, Validation.Length(3, 50)),
	)
}

func (r *Rules) ValidateUpdate(model interface{}) error {

	assessmentModel, err := securityScoreCard.InterfaceToModel(model)
	if err != nil {
		return err
	}
	return Validation.ValidateStruct(assessmentModel,
		Validation.Field(&assessmentModel.DomainName, Validation.Required, is.UUIDv4),
		Validation.Field(&assessmentModel.UpdatedBy, Validation.Required, Validation.Length(3, 50)),
	)
}

func (r *Rules) ConvertIoReaderToStruct(data io.Reader, model interface{}) (interface{}, error) {
	if data == nil {
		return nil, errors.New("body is invalid")
	}
	return model, json.NewDecoder(data).Decode(model)
}
