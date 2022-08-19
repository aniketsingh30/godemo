package controldatainsert

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	EntityControlData "riscvue.com/pkg/entities/controldata"
	"riscvue.com/pkg/repository/adapter"
	ServiceInterface "riscvue.com/pkg/services"
	"riscvue.com/pkg/utils/logger"
)

const CLASSS_NAME = "ControlDataService"

type ControlDataService struct {
	ServiceInterface.ControlDataInterface
	Repository adapter.Interface
	//Req        *http.Request
}

const CONTROLDATAID = "ControlDataId"

func NewControlDataService(repository adapter.Interface) ServiceInterface.ControlDataInterface {
	return &ControlDataService{
		Repository: repository,
		//Req:        req,
	}
}

func (h *ControlDataService) FindControlsByControlDataId(controlDataId string) (entities []EntityControlData.ControlData, err error) {
	const MethodName = "FindControlsByControlDataId"
	entities = []EntityControlData.ControlData{}
	var entity EntityControlData.ControlData
	var filter expression.ConditionBuilder
	if controlDataId != "" {
		filter = expression.Name(CONTROLDATAID).Equal(expression.Value(controlDataId))

		condition, err := expression.NewBuilder().WithFilter(filter).Build()
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:expression.NewBuilder:::", err)
			return entities, err
		}

		response, err := h.Repository.FindAll(condition, entity.TableName())
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:FindAll:::", err)
			return entities, err
		}

		if response != nil {
			for _, value := range response.Items {
				entity, err := EntityControlData.ParseDynamoAtributeToStruct(value)
				if err != nil {
					logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:ParseDynamoAtributeToStruct:::", err)
					return entities, err
				}
				entities = append(entities, entity)
			}
		}
	}
	return entities, nil
}
