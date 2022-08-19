package user

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

type User struct {
	entities.Base
	UserId               string    `json:"userId"`
	FullName             string    `json:fullName"`
	CustomerId           string    `json:"customerId"`
	IsPaid               bool      `json:"isPaid"`
	IsSuperAdmin         bool      `json:"isSuperAdmin"`
	OrganizationName     string    `json:"organizationName"`
	LastAccessCustomerID string    `json:"lastAccessCustomerID"`
	LastActive           time.Time `json:"lastActive"`
	Profile              string    `json:"profile"`
	SecretCode           string    `json:"secretCode"`
}

func InterfaceToModel(data interface{}) (instance *User, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}

	return instance, json.Unmarshal(bytes, &instance)
}

func (p *User) GetFilterId() map[string]interface{} {
	return map[string]interface{}{"UserId": p.UserId}
}

func (p *User) TableName() string {
	return os.Getenv("STORAGE_USERTABLE_NAME") + "-" + os.Getenv("ENV")
}

func (p *User) Bytes() ([]byte, error) {
	return json.Marshal(p)
}

func (p *User) GetMap() map[string]interface{} {
	logger.INFO("className=UserEntitiy MethodName=GetMap start:::", nil)
	return map[string]interface{}{
		"id":                   p.ID.String(),
		"UserId":               p.UserId,
		"CustomerId":           p.CustomerId,
		"FullName":             p.FullName,
		"IsPaid":               p.IsPaid,
		"IsSuperAdmin":         p.IsSuperAdmin,
		"Profile":              p.Profile,
		"OrganizationName":     p.OrganizationName,
		"SecretCode":           p.SecretCode,
		"LastAccessCustomerID": p.LastAccessCustomerID,
		"LastActive":           p.LastActive.Format(entities.GetTimeFormat()),
		"createdAt":            p.CreatedAt.Format(entities.GetTimeFormat()),
		"updatedAt":            p.UpdatedAt.Format(entities.GetTimeFormat()),
		"createdBy":            p.UserId,
		"updatedBy":            p.UserId,
	}
}

func ParseDynamoAtributeToStruct(response map[string]*dynamodb.AttributeValue) (p User, err error) {
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

		if key == "FullName" {
			p.FullName = *value.S
		}
		if key == "Profile" {
			p.Profile = *value.S
		}
		if key == "LastActive" {
			p.LastActive, err = time.Parse(entities.GetTimeFormat(), *value.S)
		}
		if key == "IsPaid" {
			p.IsPaid = *value.BOOL
		}
		if key == "IsSuperAdmin" {
			p.IsSuperAdmin = *value.BOOL
		}
		if key == "OrganizationName" {
			p.OrganizationName = *value.S
		}
		if key == "SecretCode" {
			p.SecretCode = *value.S
		}

		if key == "LastAccessCustomerID" {
			p.LastAccessCustomerID = *value.S
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
