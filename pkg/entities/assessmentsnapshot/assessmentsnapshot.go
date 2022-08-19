package assessment

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/google/uuid"
	"riscvue.com/pkg/entities"
	EntityAssessmentRoleMapping "riscvue.com/pkg/entities/assessmentrolemapping"
	EntityAssessmentControl "riscvue.com/pkg/entities/control"
	"riscvue.com/pkg/utils/logger"
)

const CLASSS_NAME = "AssessmentEntity"

type AssessmentSnapshot struct {
	entities.Base
	AssessmentId               uuid.UUID                                               `json:"assessmentId"`
	AssessmentUserMappingRoles []EntityAssessmentRoleMapping.AssessmentUserMappingRole `json:"assessmentUsers"`
	Name                       string                                                  `json:"name"`
	SecurityFrameworkId        string                                                  `json:"securityFrameworkId"`
	CustomerId                 string                                                  `json:"customerId"`
	MaturityFrameworkId        string                                                  `json:"maturityFrameworkId"`
	Score                      float64                                                 `json:"score"`
	Progress                   string                                                  `json:"progress"`
	LastActivity               time.Time                                               `json:"lastActivity"`
	AssessmentControls         []EntityAssessmentControl.AssessmentControl             `json:"assessmentControls"`
}

func InterfaceToModel(data interface{}) (instance *AssessmentSnapshot, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}

	return instance, json.Unmarshal(bytes, &instance)
}

func (p *AssessmentSnapshot) GetFilterId() map[string]interface{} {
	return map[string]interface{}{"id": p.ID.String()}
}

func (p *AssessmentSnapshot) TableName() string {
	return os.Getenv("STORAGE_ASSESMENT_SNAPSHOT_TABLE_NAME") + "-" + os.Getenv("ENV")

}

func (p *AssessmentSnapshot) Bytes() ([]byte, error) {
	return json.Marshal(p)
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func (p *AssessmentSnapshot) GetMap() map[string]interface{} {
	logger.INFO("className=AssessmentSnapshot MethodName=GetMap start:::", nil)
	score := FloatToString(p.Score)

	return map[string]interface{}{
		"id":                         p.ID.String(),
		"AssessmentId":               p.AssessmentId.String(),
		"Name":                       p.Name,
		"SecurityFrameworkId":        p.SecurityFrameworkId,
		"MaturityFrameworkId":        p.MaturityFrameworkId,
		"Progress":                   p.Progress,
		"Score":                      score,
		"CustomerId":                 p.CustomerId,
		"AssessmentUserMappingRoles": p.AssessmentUserMappingRoles,
		"AssessmentControls":         p.AssessmentControls,

		"LastActivity": p.LastActivity.Format(entities.GetTimeFormat()),
		"createdAt":    p.CreatedAt.Format(entities.GetTimeFormat()),
		"updatedAt":    p.UpdatedAt.Format(entities.GetTimeFormat()),
		"createdBy":    p.CreatedBy,
		"updatedBy":    p.CreatedBy,
	}
}

func ParseDynamoAtributeToStruct(response map[string]*dynamodb.AttributeValue) (p AssessmentSnapshot, err error) {
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
		if key == "AssessmentId" {
			p.AssessmentId, err = uuid.Parse(*value.S)
			if p.AssessmentId == uuid.Nil {
				err = errors.New("Item not found")
			}
		}
		if key == "Name" {
			p.Name = *value.S
		}
		if key == "SecurityFrameworkId" {
			p.SecurityFrameworkId = *value.S
		}

		if key == "MaturityFrameworkId" {
			p.MaturityFrameworkId = *value.S
		}
		if key == "Score" {
			p.Score, _ = strconv.ParseFloat(*value.S, 64)
		}
		if key == "Progress" {
			p.Progress = *value.S
		}

		if key == "CustomerId" {
			p.CustomerId = *value.S
		}

		if key == "createdBy" {
			p.CreatedBy = *value.S
		}
		if key == "updatedBy" {
			p.UpdatedBy = *value.S
		}

		if key == "AssessmentUserMappingRoles" {
			if value.L != nil {
				dynamodbattribute.UnmarshalList(value.L, &p.AssessmentUserMappingRoles)

			} else {

				p.AssessmentUserMappingRoles = []EntityAssessmentRoleMapping.AssessmentUserMappingRole{}
			}

		}

		if key == "AssessmentControls" {
			if value.L != nil {
				dynamodbattribute.UnmarshalList(value.L, &p.AssessmentControls)

			} else {

				p.AssessmentControls = []EntityAssessmentControl.AssessmentControl{}
			}

		}

		if key == "LastActivity" {
			p.LastActivity, err = time.Parse(entities.GetTimeFormat(), *value.S)
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
