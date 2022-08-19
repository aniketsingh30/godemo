package controldatainsert

import (
	"fmt"
	"os"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"riscvue.com/pkg/entities/controldata"
	EntityControlData "riscvue.com/pkg/entities/controldata"
	"riscvue.com/pkg/repository/adapter"
	ServiceInterface "riscvue.com/pkg/services"
	"riscvue.com/pkg/utils/logger"
)

const CLASSS_NAME = "ControldataInsertService"
const RESULT = "result"
const STATUS_CODE = "status"
const SUCCESS_MSG = "successMessage"
const FAILLED_MSG = "failedMessage"

var ErrorMethodNotAllowed = "method Not allowed"

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}
type ControlDataInsertService struct {
	ServiceInterface.ControlDataInterface
	Repository adapter.Interface
}

const SUCCESS_RETRIVED = "Successfully Retrieved."
const SUCCESS_DELETED = "Successfully Deleted."
const SUCCESS_UPDATED = "Successfully Updated."
const SUCCESS_CREATED = "Successfully Created."

func NewControlDataInsertService(repository adapter.Interface) ServiceInterface.ControlDataInsertInterface {
	return &ControlDataInsertService{
		Repository: repository,
		//Req:        req,
		//Rules:      RulesControldata.NewRules(),
	}
}

func (h *ControlDataInsertService) CreateRecord(controlDataReq *EntityControlData.ControlDataRequest) (controlDatas []EntityControlData.ControlData, err error) {
	const MethodName = "CreateRecord"
	controlDatas = []EntityControlData.ControlData{}

	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" start:::", controlDataReq)
	allData, err := readExcel(controlDataReq)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:validate::", err)
		return controlDatas, err
	}
	for _, excelData := range allData[4:] {
		entity, err := h.CreateRecordInternal(excelData)
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:CreateRecordInterNal::", err)
			return controlDatas, err
		}
		controlDatas = append(controlDatas, entity)
	}

	return controlDatas, nil

}

func (h *ControlDataInsertService) CreateRecordFromOtherEnv(env string) (controlDatas []EntityControlData.ControlData, err error) {
	const MethodName = "CreateRecordFromOtherEnv"
	controlDatas = []EntityControlData.ControlData{}
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" start:::", env)
	env = "dev"
	filter := expression.Name("name").NotEqual(expression.Value(""))
	condition, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return controlDatas, err
	}
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" tableName:::", os.Getenv("STORAGE_CONTROLDATATABLE_NAME")+"-"+env)
	response, err := h.Repository.FindAll(condition, os.Getenv("STORAGE_CONTROLDATATABLE_NAME")+"-"+env)
	if err != nil {
		return controlDatas, err
	}
	if response != nil {
		for i, value := range response.Items {
			entity, err := controldata.ParseDynamoAtributeToStruct(value)
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" ent:::", i)
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" ent:::", entity.TableName())
			if err != nil {
				return controlDatas, err
			}
			resp, err := h.Repository.CreateOrUpdate(entity.GetMap(), entity.TableName())

			if err != nil {
				logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" start:::", err)
				return controlDatas, err
			}
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" success insert:::", i)
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" start:::", resp)

			controlDatas = append(controlDatas, entity)
		}
	}

	return controlDatas, nil

}
func (h *ControlDataInsertService) CreateRecordInternal(excelData ServiceInterface.AllData) (entity EntityControlData.ControlData, err error) {
	const MethodName = "CreateRecordInternal"
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" start:::", excelData)
	var entityBody EntityControlData.ControlData
	entityBody.ControlNumber = excelData.ControlNumber
	entityBody.ControlFamily = excelData.ControlFamily
	entityBody.MappedControl = excelData.MappedControl
	entityBody.ShortDescription = excelData.ShortDescription
	entityBody.Requirement = excelData.Requirement

	setDefaultValues(&entityBody, entityBody.ControlDataId)
	entityBody.UpdatedBy = entityBody.CreatedBy
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" entityBody:::", entityBody)
	resp, err := h.Repository.CreateOrUpdate(entityBody.GetMap(), entityBody.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:CreateOrUpdate::", err)
		return entityBody, err
	}
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" response:::", resp)

	return entityBody, nil

}

func setDefaultValues(assesment *EntityControlData.ControlData, controlDataId string) {
	assesment.UpdatedAt = time.Now()
	if controlDataId == "" {
		assesment.CreatedAt = time.Now()
	}
	assesment.ControlDataId = controlDataId

}

func readExcel(controlDataReq *EntityControlData.ControlDataRequest) (mapping []ServiceInterface.AllData, err error) {
	mapping = []ServiceInterface.AllData{}
	fileName := controlDataReq.FilePath

	xlsx, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println(err)
		return mapping, err
	}

	rows := xlsx.GetRows(controlDataReq.SheetName)
	for _, row := range rows[4:] {
		mapping = append(mapping, ServiceInterface.AllData{ControlNumber: row[0], ControlFamily: row[1], ShortDescription: row[2], MappedControl: row[3], Requirement: row[4]})
	}

	// m := map[string]interface{}{
	//	"Data": mapping,
	//}

	//fmt.Println(mapping)
	return mapping, nil
}
