package controldata

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

type ControlDataRequest struct {
	ControlDataId string `json:"controlDataId"`
	FilePath      string `json:"filePath"`
	SheetName     string `json:"sheetName"`
	CreatedBy     string `json:"createdBy"`
}

type ControlData struct {
	entities.Base
	ControlDataId    string `json:"controlDataId"`
	ControlNumber    string `json:"controlNumber"`
	ControlFamily    string `json:"controlFamily"`
	ShortDescription string `json:"shortDescription"`
	MappedControl    string `json:"MappedControl"`
	Requirement      string `json:"requirement"`
}

func InterfaceToModel(data interface{}) (instance *ControlData, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}

	return instance, json.Unmarshal(bytes, &instance)
}

func InterfaceToModelReq(data interface{}) (instance *ControlDataRequest, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}

	return instance, json.Unmarshal(bytes, &instance)
}
func (p *ControlData) GetFilterId() map[string]interface{} {
	return map[string]interface{}{"controlDataId": p.ControlDataId}
}

func (p *ControlData) TableName() string {
	return os.Getenv("STORAGE_CONTROLDATATABLE_NAME") + "-" + os.Getenv("ENV")
}

func (p *ControlData) Bytes() ([]byte, error) {
	return json.Marshal(p)
}

func (p *ControlData) GetMap() map[string]interface{} {
	logger.INFO("className=Entitiy MethodName=GetMap start:::", nil)
	return map[string]interface{}{
		"id":            p.ID.String(),
		"ControlDataId": p.ControlDataId,
		//"Questions": p.Questions,
		"ControlNumber":    p.ControlNumber,
		"ControlFamily":    p.ControlFamily,
		"ShortDescription": p.ShortDescription,
		"MappedControl":    p.MappedControl,
		"Requirement":      p.Requirement,
		"createdAt":        p.CreatedAt.Format(entities.GetTimeFormat()),
		"updatedAt":        p.UpdatedAt.Format(entities.GetTimeFormat()),
		"createdBy":        p.CreatedBy,
		"updatedBy":        p.CreatedBy,
	}
}

func ParseDynamoAtributeToStruct(response map[string]*dynamodb.AttributeValue) (p ControlData, err error) {
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

		if key == "ControlDataId" {
			p.ControlDataId = *value.S
		}
		// if key == "Questions" {

		// 		if value.L != nil {
		// 			dynamodbattribute.UnmarshalList(value.L, &p.Questions)

		// 		} else {

		// 			p.Questions = []string{}
		// 		}

		// }
		if key == "ControlNumber" {
			p.ControlNumber = *value.S
		}
		if key == "ControlFamily" {
			p.ControlFamily = *value.S
		}
		if key == "ShortDescription" {
			p.ShortDescription = *value.S
		}
		if key == "MappedControl" {
			p.MappedControl = *value.S
		}

		if key == "Requirement" {
			p.Requirement = *value.S
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
