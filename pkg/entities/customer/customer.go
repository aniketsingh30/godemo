package customer

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"riscvue.com/pkg/entities"
	"riscvue.com/pkg/utils/logger"
)

type Customer struct {
	CustomerId   string        `json:"customerId"`
	CustomerName string        `json:"customerName"`
	OwnerId      string        `json:"ownerId"`
	Address      string        `json:"address"`
	State        string        `json:"state"`
	ZipCode      string        `json:"zipCode"`
	Country      string        `json:"country"`
	NetWork      NetWork       `json:"network"`
	Vendors      []string      `json:"vendors"`
	Integration  []Integration `json:"integrations"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	CreatedBy    string        `json:"createdBy"`
	UpdatedBy    string        `json:"updatedBy"`
}

type CustomerUpdate struct {
	FieldName    string        `json:"fieldName"`
	CustomerId   string        `json:"customerId"`
	CustomerName string        `json:"customerName"`
	OwnerId      string        `json:"ownerId"`
	Address      string        `json:"address"`
	State        string        `json:"state"`
	ZipCode      string        `json:"zipCode"`
	Country      string        `json:"country"`
	NetWork      NetWork       `json:"network"`
	Vendors      []string      `json:"vendors"`
	Integration  []Integration `json:"integrations"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	CreatedBy    string        `json:"createdBy"`
	UpdatedBy    string        `json:"updatedBy"`
}

type NetWork struct {
	DNSDomain string `json:"dnsDomain"`
	CIDR      string `json:"cidr"`
	Website   string `json:"website"`
}
type Integration struct {
	Name  string `json:"name"`
	Value bool   `json:"value"`
}

func InterfaceToModel(data interface{}) (instance *Customer, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}

	return instance, json.Unmarshal(bytes, &instance)
}

func InterfaceToModelUpdate(data interface{}) (instance *CustomerUpdate, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}

	return instance, json.Unmarshal(bytes, &instance)
}

func (p *Customer) GetFilterId() map[string]interface{} {
	return map[string]interface{}{"CustomerId": p.CustomerId}
}

func (p *Customer) TableName() string {
	return os.Getenv("STORAGE_CUSTOMERTABLE_NAME") + "-" + os.Getenv("ENV")
}

func (p *Customer) Bytes() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Customer) GetMap() map[string]interface{} {
	logger.INFO("className=CustomerEntitiy MethodName=GetMap start:::", nil)
	return map[string]interface{}{
		"CustomerId":   p.CustomerId,
		"CustomerName": p.CustomerName,
		"OwnerId":      p.OwnerId,
		"Address":      p.Address,
		"State":        p.State,
		"ZipCode":      p.ZipCode,
		"Country":      p.Country,
		"NetWork":      p.NetWork,
		"Vendors":      p.Vendors,
		"Integration":  p.Integration,
		"createdAt":    p.CreatedAt.Format(entities.GetTimeFormat()),
		"updatedAt":    p.UpdatedAt.Format(entities.GetTimeFormat()),
		"createdBy":    p.OwnerId,
		"updatedBy":    p.OwnerId,
	}
}

func ParseDynamoAtributeToStruct(response map[string]*dynamodb.AttributeValue) (p Customer, err error) {
	if response == nil || (response != nil && len(response) == 0) {
		return p, errors.New("Item not found")
	}
	for key, value := range response {

		if key == "CustomerId" {
			p.CustomerId = *value.S
		}
		if key == "CustomerId" {
			p.CustomerId = *value.S
		}

		if key == "CustomerName" {
			p.CustomerName = *value.S
		}
		if key == "OwnerId" {
			p.OwnerId = *value.S
		}

		if key == "Address" {
			p.Address = *value.S
		}
		if key == "State" {
			p.State = *value.S
		}

		if key == "ZipCode" {
			p.ZipCode = *value.S
		}

		if key == "Country" {
			p.Country = *value.S
		}

		if key == "NetWork" {
			if value.M != nil {
				dynamodbattribute.UnmarshalMap(value.M, &p.NetWork)

			}
		}
		if key == "Vendors" {
			if value.L != nil {
				dynamodbattribute.UnmarshalList(value.L, &p.Vendors)
			}
		}

		if key == "Integration" {
			if value.L != nil {
				dynamodbattribute.UnmarshalList(value.L, &p.Integration)

			}
		}

		if key == "createdBy" {
			p.CreatedBy = *value.S
		}
		if key == "updatedBy" {
			p.UpdatedBy = *value.S
		}
		if key == "createdAt" {
			p.CreatedAt, err = time.Parse(entities.GetTimeFormat(), *value.S)
		}
		if key == "updatedAt" {
			p.UpdatedAt, err = time.Parse(entities.GetTimeFormat(), *value.S)
		}
		if err != nil {
			return p, err
		}
	}

	return p, nil
}
