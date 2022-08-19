package usermapping

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"riscvue.com/pkg/entities"
	"riscvue.com/pkg/utils/logger"
)

type UserMapping struct {
	entities.Base
	UserId           string `json:"userId"`
	CustomerId       string `json:"customerId"`
	IsPaid           string `json:"isPaid"`
	FullName         string `json:"fullName"`
	IsSuperAdmin     string `json:"isSuperAdmin"`
	IsActive         string `json:"isActive"`
	OrganizationName string `json:"organizationName"`
}

func InterfaceToModel(data interface{}) (instance *UserMapping, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}

	return instance, json.Unmarshal(bytes, &instance)
}

func (p *UserMapping) GetFilterId() map[string]interface{} {
	return map[string]interface{}{"UserId": p.UserId}
}

func (p *UserMapping) TableName() string {
	return os.Getenv("STORAGE_USER_MAPPING_TABLE_NAME") + "-" + os.Getenv("ENV")
}

func (p *UserMapping) Bytes() ([]byte, error) {
	return json.Marshal(p)
}

func (p *UserMapping) GetMap() map[string]interface{} {
	logger.INFO("className=UserMappingEntitiy MethodName=GetMap start:::", nil)
	return map[string]interface{}{
		"id":               p.ID.String(),
		"UserId":           p.UserId,
		"CustomerId":       p.CustomerId,
		"createdBy":        p.UserId,
		"updatedBy":        p.UserId,
		"IsPaid":           p.IsPaid,
		"FullName":         p.FullName,
		"IsSuperAdmin":     p.IsSuperAdmin,
		"IsActive":         p.IsActive,
		"OrganizationName": p.OrganizationName,
		"createdAt":        p.CreatedAt.Format(entities.GetTimeFormat()),
		"updatedAt":        p.UpdatedAt.Format(entities.GetTimeFormat()),
	}
}

func ParseDynamoAtributeToStruct(response map[string]*dynamodb.AttributeValue) (p UserMapping, err error) {
	if response == nil || (response != nil && len(response) == 0) {
		return p, errors.New("Item not found")
	}
	for key, value := range response {
		if key == "id" {
			p.ID, err = uuid.Parse(*value.S)
			if p.ID == uuid.Nil {
				err = errors.New("Item not found")
			}
		}
		if key == "UserId" {
			p.UserId = *value.S
		}

		if key == "CustomerId" {
			p.CustomerId = *value.S
		}

		if key == "IsPaid" {
			p.IsPaid = *value.S
		}
		if key == "IsActive" {
			p.IsActive = *value.S
		}
		if key == "FullName" {
			p.FullName = *value.S
		}
		if key == "IsSuperAdmin" {
			p.IsSuperAdmin = *value.S
		}
		if key == "OrganizationName" {
			p.OrganizationName = *value.S
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
