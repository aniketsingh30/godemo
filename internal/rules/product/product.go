package product

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	Validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"riscvue.com/internal/entities"
	"riscvue.com/internal/entities/product"
)

type Rules struct{}

func NewRules() *Rules {
	return &Rules{}
}

func (r *Rules) ConvertIoReaderToStruct(data io.Reader, model interface{}) (interface{}, error) {
	if data == nil {
		return nil, errors.New("body is invalid")
	}
	return model, json.NewDecoder(data).Decode(model)
}

func (r *Rules) Migrate(connection *dynamodb.DynamoDB) error {
	r.MaturityFrameworkTable(connection)
	r.SecurityControlFrameworkTable(connection)
	r.ControlDataTable(connection)
	r.SecurityScoreCard(connection)
	r.QualysReportTable(connection)
	r.QualysReportScanTable(connection)
	r.QualysReportAssetTable(connection)
	r.QualysScanInfoTable(connection)
	r.QualysReportScanAuditTable(connection)
	r.QualysScanInfoAuditTable(connection)
	return r.ControlDataMappingTable(connection)
}

func (r *Rules) MaturityFrameworkTable(connection *dynamodb.DynamoDB) error {
	tableName := os.Getenv("STORAGE_MATURITYFRAMEWORKTABLE_NAME") + "-" + os.Getenv("ENV")

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("MaturityFWName"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("MaturityIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{

					{
						AttributeName: aws.String("MaturityFWName"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},

				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
			},
		},

		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
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

func (r *Rules) SecurityControlFrameworkTable(connection *dynamodb.DynamoDB) error {
	tableName := os.Getenv("STORAGE_SECURITYCONTROLFWTABLE_NAME") + "-" + os.Getenv("ENV")

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("SecurityControlName"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("SecurityIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{

					{
						AttributeName: aws.String("SecurityControlName"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},

				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
			},
		},

		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
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

func (r *Rules) QualysReportTable(connection *dynamodb.DynamoDB) error {
	tableName := os.Getenv("STORAGE_QUALYS_REPORT_TABLE") + "-" + os.Getenv("ENV")

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("CustomerId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("CustomerId"),
				KeyType:       aws.String("HASH"),
			},
		},

		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
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

func (r *Rules) QualysReportHistoryTable(connection *dynamodb.DynamoDB) error {
	tableName := os.Getenv("STORAGE_QUALYS_REPORT_AUDIT_TABLE") + "-" + os.Getenv("ENV")

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("CustomerId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("CustomerIdIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{

					{
						AttributeName: aws.String("CustomerId"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},

				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
			},
		},

		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
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
func (r *Rules) QualysReportScanTable(connection *dynamodb.DynamoDB) error {
	tableName := os.Getenv("STORAGE_QUALYS_SCAN_REPORT_TABLE") + "-" + os.Getenv("ENV")

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("CustomerId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("CustomerId"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
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

func (r *Rules) QualysReportScanAuditTable(connection *dynamodb.DynamoDB) error {
	tableName := os.Getenv("STORAGE_QUALYS_SCAN_REPORT_AUDIT_TABLE") + "-" + os.Getenv("ENV")

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("CustomerId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("CustomerIdIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{

					{
						AttributeName: aws.String("CustomerId"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},

				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
			},
		},

		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
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
func (r *Rules) QualysScanInfoTable(connection *dynamodb.DynamoDB) error {
	tableName := os.Getenv("STORAGE_QUALYS_SCAN_INFO_TABLE") + "-" + os.Getenv("ENV")

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("CustomerId"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Status"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("CustomerId"),
				KeyType:       aws.String("HASH"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("StatusIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{

					{
						AttributeName: aws.String("Status"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},

				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
			},
		},

		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
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

func (r *Rules) QualysScanInfoAuditTable(connection *dynamodb.DynamoDB) error {
	tableName := os.Getenv("STORAGE_QUALYS_SCAN_INFO_AUDIT_TABLE") + "-" + os.Getenv("ENV")

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("CustomerId"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Status"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("CustomerId"),
				KeyType:       aws.String("HASH"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("StatusIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{

					{
						AttributeName: aws.String("Status"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},

				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
			},
		},

		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
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
func (r *Rules) QualysReportAssetTable(connection *dynamodb.DynamoDB) error {
	tableName := os.Getenv("STORAGE_QUALYS_ASSET_GROUP_TABLE") + "-" + os.Getenv("ENV")

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{

			{
				AttributeName: aws.String("CustomerId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("CustomerId"),
				KeyType:       aws.String("HASH"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("CustomerIdIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{

					{
						AttributeName: aws.String("CustomerId"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},

				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
			},
		},

		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
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

func (r *Rules) ControlDataTable(connection *dynamodb.DynamoDB) error {
	tableName := os.Getenv("STORAGE_CONTROLDATATABLE_NAME") + "-" + os.Getenv("ENV")

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("ControlDataId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("ControlDataIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{

					{
						AttributeName: aws.String("ControlDataId"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},

				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
			},
		},

		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
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

func (r *Rules) ControlDataMappingTable(connection *dynamodb.DynamoDB) error {
	tableName := os.Getenv("STORAGE_CONTROLMAPPINGTABLE_NAME") + "-" + os.Getenv("ENV")

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("MappingKey"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("MappingIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{

					{
						AttributeName: aws.String("MappingKey"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},

				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
			},
		},

		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
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

func (r *Rules) SecurityScoreCard(connection *dynamodb.DynamoDB) error {
	tableName := os.Getenv("STORAGE_SECURITYSCORECARD_TABLE_NAME") + "-" + os.Getenv("ENV")

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{

			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("DomainName"),
				AttributeType: aws.String("S"),
			},

			{
				AttributeName: aws.String("CustomerId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("DomainNameIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{

					{
						AttributeName: aws.String("DomainName"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},

				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
			},
			{
				IndexName: aws.String("CustomerIdIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{

					{
						AttributeName: aws.String("CustomerId"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},

				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
			},
		},

		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
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
func (r *Rules) createTable(connection *dynamodb.DynamoDB) error {
	table := &product.Product{}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("_id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("_id"),
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
	return product.Product{
		Base: entities.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: uuid.New().String(),
	}
}

func (r *Rules) Validate(model interface{}) error {
	productModel, err := product.InterfaceToModel(model)
	if err != nil {
		return err
	}

	return Validation.ValidateStruct(productModel,
		Validation.Field(&productModel.ID, Validation.Required, is.UUIDv4),
		Validation.Field(&productModel.Name, Validation.Required, Validation.Length(3, 50)),
	)
}
