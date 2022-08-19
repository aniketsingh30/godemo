package usermapping

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	Validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"riscvue.com/pkg/entities"
	"riscvue.com/pkg/entities/usermapping"
)

type Rules struct{}

func NewRules() *Rules {
	return &Rules{}
}

func (r *Rules) Migrate(connection *dynamodb.DynamoDB) error {
	return r.createTable(connection)
}

func (r *Rules) createTable(connection *dynamodb.DynamoDB) error {
	table := &usermapping.UserMapping{}

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
	return usermapping.UserMapping{
		Base: entities.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func (r *Rules) Validate(model interface{}) error {

	userMappingModel, err := usermapping.InterfaceToModel(model)
	if err != nil {
		return err
	}
	return Validation.ValidateStruct(userMappingModel,
		Validation.Field(&userMappingModel.ID, Validation.Required, is.UUIDv4),
		//Validation.Field(&userMappingModel.Username, Validation.Required, Validation.Length(3, 50)),
	)
}
