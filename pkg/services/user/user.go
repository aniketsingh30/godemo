package user

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	EntityUser "riscvue.com/pkg/entities/user"
	"riscvue.com/pkg/repository/adapter"

	ServiceInterface "riscvue.com/pkg/services"
	"riscvue.com/pkg/utils/logger"
)

const CLASSS_NAME = "UserService"

type UserService struct {
	ServiceInterface.OrganizationInterface
	Repository adapter.Interface
	//Req        *http.Request
}

const SUCCESS_RETRIVED = "Successfully Retrieved."

func NewUserService(repository adapter.Interface) ServiceInterface.UserInterface {
	return &UserService{
		Repository: repository,
		//	Req:        req,
	}
}

func (h *UserService) GetRecordByUserId(userId string) (user EntityUser.User, err error) {

	const MethodName = "GetRecordByUserId"

	var entityBody EntityUser.User
	entityBody.UserId = userId
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" TableName::"+entityBody.TableName()+" userId:::", userId)

	resp, err := h.Repository.FindOne(entityBody.GetFilterId(), entityBody.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error FindOne:::", err)
		return entityBody, err
	}
	data, err := EntityUser.ParseDynamoAtributeToStruct(resp.Item)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:ParseDynamoAtributeToStruct:::", err)
		return data, err
	}
	return data, nil

}

func (h *UserService) GetAllRecords() (entities []EntityUser.User, err error) {
	const MethodName = "GetAllRecords"
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" id:::", nil)
	entities = []EntityUser.User{}
	var entity EntityUser.User

	filter := expression.Name("CustomerId").NotEqual(expression.Value(""))
	condition, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:NewBuilder::", err)
		return entities, err
	}

	resp, err := h.Repository.FindAll(condition, entity.TableName())
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" resp::", resp)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName, err)
		return entities, err
	}

	if resp != nil {
		for _, value := range resp.Items {
			entity, err := EntityUser.ParseDynamoAtributeToStruct(value)
			if err != nil {
				logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:ParseDynamoAtributeToStruct::", err)
				return entities, err
			}
			entities = append(entities, entity)
		}
	}
	return entities, err
}
