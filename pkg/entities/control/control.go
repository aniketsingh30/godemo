package control

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

	"riscvue.com/pkg/utils/logger"
)

type HistorySlice []History

func (p HistorySlice) Len() int {
	return len(p)
}

func (p HistorySlice) Less(i, j int) bool {
	UpdatedAtJ, _ := time.Parse(entities.GetTimeFormat(), p[j].UpdatedAt)
	UpdatedAtI, _ := time.Parse(entities.GetTimeFormat(), p[i].UpdatedAt)
	return UpdatedAtJ.Before(UpdatedAtI)
}

func (p HistorySlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type AttachmentSlice []Attachment

func (p AttachmentSlice) Len() int {
	return len(p)
}

func (p AttachmentSlice) Less(i, j int) bool {
	UpdatedAtJ, _ := time.Parse(entities.GetTimeFormat(), p[j].UpdatedAt)
	UpdatedAtI, _ := time.Parse(entities.GetTimeFormat(), p[i].UpdatedAt)
	return UpdatedAtJ.Before(UpdatedAtI)
}

func (p AttachmentSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type Comment struct {
	entities.Base
	ControlNumber string `json:controlNumber"`
	Description   string `json:"description"`
}

type Question struct {
	QuestionNumber string  `json:"questionNumber"`
	Description    string  `json:"description"`
	IsEnable       bool    `json:"isEnable"`
	CurrentScore   float64 `json:"currentScore"`
	TargetScore    float64 `json:"targetScore"`
	//Answers  []Answer     `json:"answers"`
}
type Attachment struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
	CreatedBy   string    `json:"createdBy"`
	UpdatedBy   string    `json:"updatedBy"`
	Description string    `json:"description"`
	Path        string    `json:"path"`
	IsFile      bool      `json:"isFile"`
}

type History struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
	CreatedBy string    `json:"createdBy"`
	UpdatedBy string    `json:"updatedBy"`

	Description string `json:"description"`
}
type Map struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}
type AssessmentControl struct {
	entities.Base
	AssessmentId            uuid.UUID    `json:"assessmentId"`
	ControlNumber           string       `json:"controlNumber"`
	ControlFamily           string       `json:"controlFamily"`
	ShortDescription        string       `json:"shortDescription"`
	MappedControl           string       `json:"MappedControl"`
	Requirement             string       `json:"requirement"`
	IsEnable                bool         `json:"isEnable"`
	Status                  Map          `json:"status"`
	CurrentScore            float64      `json:"currentScore"`
	TargetScore             float64      `json:"targetScore"`
	DueDate                 time.Time    `json:"dueDate"`
	AssignedTo              string       `json:"assignedTo"`
	Likelihood              Map          `json:"likelihood"`
	Impact                  Map          `json:"impact"`
	InherentRisk            float64      `json:"inherentRisk"`
	ResidualRisk            float64      `json:"residualRisk"`
	Cost                    int64        `json:"cost"`
	ControlRisk             string       `json:"controlRisk"`
	MaturityLevel           string       `json:"maturityLevel"`
	ControlBusinessImpact   string       `json:"controlBusinessImpact"`
	ControlSecurityPriority string       `json:"controlSecurityPriority"`
	Comment                 Comment      `json:"comment"`
	Questions               []Question   `json:"questions"`
	Attachments             []Attachment `json:"attachments"`

	HistoryData []History `json:"historyData"`
	CreatedBy   string    `json:"createdBy"`
	UpdatedBy   string    `json:"updatedBy"`
}

// type Control struct {
// 	entities.Base
// 	AssessmentControls []AssessmentControl `json:"assessmentControls"`
// }

func InterfaceToModel(data interface{}) (instance *AssessmentControl, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}

	return instance, json.Unmarshal(bytes, &instance)
}

func (p *AssessmentControl) TableName() string {
	return os.Getenv("STORAGE_CONTROL_TABLE_NAME") + "-" + os.Getenv("ENV")
}

func (p *AssessmentControl) GetFilterId() map[string]interface{} {
	return map[string]interface{}{"id": p.ID.String()}
}
func (p *AssessmentControl) Bytes() ([]byte, error) {
	return json.Marshal(p)
}

func (p *AssessmentControl) GetMap() map[string]interface{} {
	logger.INFO("className=ControlEntitiy MethodName=GetMap start:::", nil)
	return map[string]interface{}{
		"id":                      p.ID.String(),
		"AssessmentId":            p.AssessmentId.String(),
		"ControlNumber":           p.ControlNumber,
		"ControlFamily":           p.ControlFamily,
		"ShortDescription":        p.ShortDescription,
		"MappedControl":           p.MappedControl,
		"Requirement":             p.Requirement,
		"IsEnable":                p.IsEnable,
		"Status":                  p.Status,
		"CurrentScore":            p.CurrentScore,
		"TargetScore":             p.TargetScore,
		"DueDate":                 p.DueDate.Format(entities.GetTimeFormat()),
		"AssignedTo":              p.AssignedTo,
		"Likelihood":              p.Likelihood,
		"Impact":                  p.Impact,
		"InherentRisk":            p.InherentRisk,
		"ResidualRisk":            p.ResidualRisk,
		"Cost":                    p.Cost,
		"ControlRisk":             p.ControlRisk,
		"MaturityLevel":           p.MaturityLevel,
		"ControlBusinessImpact":   p.ControlBusinessImpact,
		"ControlSecurityPriority": p.ControlSecurityPriority,
		"Comment":                 p.Comment,
		"Questions":               p.Questions,
		"Attachments":             p.Attachments,
		"HistoryData":             p.HistoryData,
		"createdAt":               p.CreatedAt.Format(entities.GetTimeFormat()),
		"updatedAt":               p.UpdatedAt.Format(entities.GetTimeFormat()),
		"createdBy":               p.CreatedBy,
		"updatedBy":               p.UpdatedBy,
	}
}

func ParseDynamoAtributeToStruct(response map[string]*dynamodb.AttributeValue) (p AssessmentControl, err error) {
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

		if key == "ControlNumber" {
			p.ControlNumber = *value.S
		}
		if key == "createdBy" {
			p.CreatedBy = *value.S
		}
		if key == "updatedBy" {
			p.UpdatedBy = *value.S
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
		if key == "IsEnable" {
			p.IsEnable = *value.BOOL
		}
		if key == "Status" {
			if value.M != nil {
				dynamodbattribute.UnmarshalMap(value.M, &p.Status)

			}
		}
		if key == "CurrentScore" {
			p.CurrentScore, _ = strconv.ParseFloat(*value.N, 64)

		}
		if key == "TargetScore" {
			p.TargetScore, _ = strconv.ParseFloat(*value.N, 64)
		}
		if key == "DueDate" {
			p.DueDate, err = time.Parse(entities.GetTimeFormat(), *value.S)
		}
		if key == "AssignedTo" {
			p.AssignedTo = *value.S
		}

		if key == "Likelihood" {
			if value.M != nil {
				dynamodbattribute.UnmarshalMap(value.M, &p.Likelihood)

			}
		}
		if key == "Impact" {
			if value.M != nil {
				dynamodbattribute.UnmarshalMap(value.M, &p.Impact)

			}
		}

		if key == "InherentRisk" {
			p.InherentRisk, _ = strconv.ParseFloat(*value.N, 64)
		}
		if key == "ResidualRisk" {
			p.ResidualRisk, _ = strconv.ParseFloat(*value.N, 64)
		}

		if key == "Cost" {
			p.Cost, _ = strconv.ParseInt(*value.N, 0, 64)
		}
		if key == "ControlRisk" {
			p.ControlRisk = *value.S
		}
		if key == "ControlBusinessImpact" {
			p.ControlBusinessImpact = *value.S
		}
		if key == "ControlSecurityPriority" {
			p.ControlSecurityPriority = *value.S
		}

		if key == "MaturityLevel" {
			p.MaturityLevel = *value.S
		}
		if key == "Comment" {
			if value.M != nil {
				dynamodbattribute.UnmarshalMap(value.M, &p.Comment)

			}
		}
		if key == "Questions" {
			if value.L != nil {
				dynamodbattribute.UnmarshalList(value.L, &p.Questions)

			} else {

				p.Questions = []Question{}
			}
		}
		if key == "Attachments" {
			if value.L != nil {
				dynamodbattribute.UnmarshalList(value.L, &p.Attachments)

			} else {

				p.Attachments = []Attachment{}
			}
		}
		if key == "HistoryData" {
			if value.L != nil {
				dynamodbattribute.UnmarshalList(value.L, &p.HistoryData)

			} else {

				p.HistoryData = []History{}
			}

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
