package securityScoreCard

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"riscvue.com/pkg/entities"
	"riscvue.com/pkg/utils/logger"
)

const CLASSS_NAME = "SecurityScoreCard"

type MainResponse struct {
	entities.Base
	DomainName    string        `json:"domainName"`
	ScoreResponse ScoreResponse `json:"scoreResponse"`
	Factors       Factors       `json:"factors"`
	CustomerId    string        `json:"customerId"`
}

type Factors struct {
	Entries []Entries `json:"entries"`
	Total   int64     `json:"total"`
}
type Entries struct {
	Name         string         `json:"name"`
	Score        int64          `json:"score"`
	Grade        string         `json:"grade"`
	Grade_url    string         `json:"grade_url"`
	IssueSummary []IssueSummary `json:"issue_summary"`
}
type IssueSummary struct {
	Type               string `json:"type"`
	Count              int64  `json:"count"`
	Severity           string `json:"severity"`
	Total_Score_Impact string `json:"total_score_impact"`
	Detail_Url         string `json:"detail_url"`
}
type ScoreResponse struct {
	Name                   string    `json:"name"`
	Description            string    `json:"description"`
	Domain                 string    `json:"domain"`
	Grade_url              string    `json:"grade_url"`
	Industry               string    `json:"industry"`
	Size                   string    `json:"size"`
	Score                  int64     `json:"score"`
	Grade                  string    `json:"grade"`
	Last30day_score_change int64     `json:"last30day_score_change"`
	Created_at             time.Time `json:"created_at"`
	Disputed               bool      `json:"disputed"`
}

func InterfaceToModel(data interface{}) (instance *MainResponse, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}

	return instance, json.Unmarshal(bytes, &instance)
}

func (p *MainResponse) GetFilterId() map[string]interface{} {
	return map[string]interface{}{"DomainName": p.DomainName}
}

func (p *MainResponse) TableName() string {
	return os.Getenv("STORAGE_SECURITYSCORECARD_TABLE_NAME") + "-" + os.Getenv("ENV")

}

func (p *MainResponse) Bytes() ([]byte, error) {
	return json.Marshal(p)
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func (p *MainResponse) GetMap() map[string]interface{} {
	logger.INFO("className=SecurityScoreCard MethodName=GetMap start:::", nil)

	return map[string]interface{}{
		"id":            p.ID.String(),
		"DomainName":    p.DomainName,
		"ScoreResponse": p.ScoreResponse,
		"CustomerId":    p.CustomerId,
		"Factors":       p.Factors,
		"createdAt":     p.CreatedAt.Format(entities.GetTimeFormat()),
		"updatedAt":     p.UpdatedAt.Format(entities.GetTimeFormat()),
		"createdBy":     p.CreatedBy,
		"updatedBy":     p.CreatedBy,
	}
}

func ParseDynamoAtributeToStruct(response map[string]*dynamodb.AttributeValue) (p MainResponse, err error) {
	if response == nil || (response != nil && len(response) == 0) {
		return p, errors.New("Item not found")
	}
	for key, value := range response {
		if key == "DomainName" {
			p.DomainName = *value.S
		}
		if p.DomainName == "" {
			err = errors.New("Item not found")
		}

		if key == "ScoreResponse" {
			if value.M != nil {
				dynamodbattribute.UnmarshalMap(value.M, &p.ScoreResponse)

			}
		}
		if key == "Factors" {
			if value.M != nil {
				dynamodbattribute.UnmarshalMap(value.M, &p.Factors)

			}
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
