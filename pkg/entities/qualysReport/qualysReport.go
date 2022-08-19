package qualysReport

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"

	"riscvue.com/pkg/entities"
	"riscvue.com/pkg/utils/logger"
)

const CLASSS_NAME = "SecurityScoreCard"

type AddIPRequest struct {
	EnableVm string `json:enableVm`
	IPs      string `json:ips`
}

type AddEditAssestGroupRequest struct {
	ScanTitle string `json:scanTitle`
	IPs       string `json:ips`
	ScanId    string `json:scanId`
}

type LaunchScanRequest struct {
	IP          string `json:ip`
	ScanTitle   string `json:scanTitle`
	AssetGroups string `json:assetGroups`
	ClientName  string `json:"clientName"`
	OptionId    string `json:"optionId"`
}

type QualysResponse struct {
	entities.Base
	SIMPLE_RETURN struct {
		RESPONSE struct {
			DATETIME  time.Time   `json:"dateTime"`
			CODE      string      `json:"code"`
			TEXT      string      `json:"text"`
			ITEM_LIST interface{} `json:"ITEM_LIST"`
		}
	}
	CustomerId string `json:"customerId"`
	TxnType    string `json:txnType`
}

type AssetResponse struct {
	entities.Base
	SIMPLE_RETURN struct {
		RESPONSE struct {
			DATETIME  time.Time `json:"dateTime"`
			TEXT      string    `json:"text"`
			ITEM_LIST struct {
				ITEM struct {
					KEY   string `json:"key"`
					VALUE string `json:"value"`
				}
			}
		}
	}
	CustomerId string `json:"customerId"`
	TxnType    string `json:txnType`
}
type RESPONSE struct {
	DATETIME  time.Time `json:"dateTime"`
	CODE      string    `json:"code"`
	TEXT      string    `json:"text"`
	ITEM_LIST struct {
		ITEM struct {
			KEY   string `json:"key"`
			VALUE string `json:"value"`
		}
	}
}

func InterfaceToModel(data interface{}) (instance *QualysResponse, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}

	return instance, json.Unmarshal(bytes, &instance)
}

func (p *QualysResponse) GetFilterId() map[string]interface{} {
	return map[string]interface{}{"id": p.ID.String()}
}

func (p *QualysResponse) TableName() string {
	return os.Getenv("STORAGE_QUALYS_REPORT_TABLE") + "-" + os.Getenv("ENV")

}

func (p *QualysResponse) Bytes() ([]byte, error) {
	return json.Marshal(p)
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func (p *QualysResponse) GetMap() map[string]interface{} {
	logger.INFO("className=SecurityScoreCard MethodName=GetMap start:::", nil)

	return map[string]interface{}{
		"id":           p.ID.String(),
		"SimpleReturn": p.SIMPLE_RETURN,
		//"Response":     p.RESPONSE,
		"CustomerId": p.CustomerId,
		"TxnType":    p.TxnType,
		"createdAt":  p.CreatedAt.Format(entities.GetTimeFormat()),
		"updatedAt":  p.UpdatedAt.Format(entities.GetTimeFormat()),
		"createdBy":  p.CreatedBy,
		"updatedBy":  p.CreatedBy,
	}
}

func ParseDynamoAtributeToStruct(response map[string]*dynamodb.AttributeValue) (p QualysResponse, err error) {
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

		// if key == "Response" {
		// 	if value.M != nil {
		// 		dynamodbattribute.UnmarshalMap(value.M, &p.RESPONSE)

		// 	}
		// }

		if key == "CustomerId" {
			p.CustomerId = *value.S
		}
		if key == "createdBy" {
			p.CreatedBy = *value.S
		}
		if key == "updatedBy" {
			p.UpdatedBy = *value.S
		}
		if key == "TxnType" {
			p.TxnType = *value.S
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
