package invitation

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

type Invitation struct {
	entities.Base
	InvitedBy         string `json:"invitedBy"`
	InvitedTo         string `json:"invitedTo"`
	InvitedCustomerId string `json:"invitedCustomerId"`
	Status            string `json:"status"`
	IsSuperAdmin      bool   `json:"isSuperAdmin"`
	IsUserPresent     bool   `json:"isUserPresent"`
	CompanyName       string `json:"companyName"`
	FullName          string `json:"fullName"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	IsActive          bool   `json:"isActive"`
	IsPaid            bool   `json:"isPaid"`
}
type InvitationDetails struct {
	Invitation   Invitation `json:"invitation"`
	InvitedBy    string     `json:"invitedBy"`
	CustomerName string     `json:"customerName"`
}

func InterfaceToModel(data interface{}) (instance *Invitation, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}

	return instance, json.Unmarshal(bytes, &instance)
}

func (p *Invitation) GetFilterId() map[string]interface{} {
	return map[string]interface{}{"id": p.ID.String()}
}

func (p *Invitation) TableName() string {
	return os.Getenv("STORAGE_CUSTOMER_INVITATION_TABLE") + "-" + os.Getenv("ENV")
}

func (p *Invitation) Bytes() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Invitation) GetMap() map[string]interface{} {
	logger.INFO("className=Entitiy MethodName=GetMap start:::", nil)
	return map[string]interface{}{
		"id":                p.ID.String(),
		"InvitedBy":         p.InvitedBy,
		"InvitedTo":         p.InvitedTo,
		"InvitedCustomerId": p.InvitedCustomerId,
		"Status":            p.Status,
		"IsUserPresent":     p.IsUserPresent,
		"IsSuperAdmin":      p.IsSuperAdmin,
		"CompanyName":       p.CompanyName,
		"FullName":          p.FullName,
		"FirstName":         p.FirstName,
		"LastName":          p.LastName,
		"IsActive":          p.IsActive,
		"IsPaid":            p.IsPaid,
		"createdAt":         p.CreatedAt.Format(entities.GetTimeFormat()),
		"updatedAt":         p.UpdatedAt.Format(entities.GetTimeFormat()),
		"createdBy":         p.InvitedBy,
		"updatedBy":         p.InvitedBy,
	}
}

func ParseDynamoAtributeToStruct(response map[string]*dynamodb.AttributeValue) (p Invitation, err error) {
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
		if key == "InvitedCustomerId" {
			p.InvitedCustomerId = *value.S
		}
		if key == "InvitedBy" {
			p.InvitedBy = *value.S
		}

		if key == "createdBy" {
			p.CreatedBy = *value.S
		}
		if key == "updatedBy" {
			p.UpdatedBy = *value.S
		}
		if key == "InvitedTo" {
			p.InvitedTo = *value.S
		}

		if key == "IsSuperAdmin" {
			p.IsSuperAdmin = *value.BOOL
		}
		if key == "Status" {
			p.Status = *value.S
		}
		if key == "IsUserPresent" {
			p.IsUserPresent = *value.BOOL
		}
		if key == "CompanyName" {
			p.CompanyName = *value.S
		}
		if key == "FullName" {
			p.FullName = *value.S
		}
		if key == "FirstName" {
			p.FirstName = *value.S
		}
		if key == "LastName" {
			p.LastName = *value.S
		}
		if key == "IsActive" {
			p.IsActive = *value.BOOL
		}
		if key == "IsPaid" {
			p.IsPaid = *value.BOOL
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
