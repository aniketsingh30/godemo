package organization

import (
	"errors"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"

	EntityCustomer "riscvue.com/pkg/entities/customer"
	"riscvue.com/pkg/repository/adapter"
	Rules "riscvue.com/pkg/rules"
	RulesOrganization "riscvue.com/pkg/rules/organization"
	ServiceInterface "riscvue.com/pkg/services"
	"riscvue.com/pkg/utils/logger"

	ResponseUtil "riscvue.com/pkg/utils/response"
)

const CLASSS_NAME = "OrganizationService"
const RESULT = "result"
const STATUS_CODE = "status"
const SUCCESS_MSG = "successMessage"
const FAILLED_MSG = "failedMessage"

const ORGANIZATION = "Organization"
const NETWORK = "Network"
const VENDOR = "Vendor"
const INTEGRATION = "Integration"

type OrganizationService struct {
	ServiceInterface.OrganizationInterface
	Repository adapter.Interface
	Rules      Rules.Interface
}

const SUCCESS_RETRIVED = "Successfully Retrieved."
const SUCCESS_DELETED = "Successfully Deleted."
const SUCCESS_UPDATED = "Successfully Updated."
const SUCCESS_CREATED = "Successfully Created."

func NewOrganizationService(repository adapter.Interface) ServiceInterface.OrganizationInterface {
	return &OrganizationService{
		Repository: repository,
		Rules:      RulesOrganization.NewRules(),
	}
}

func (h *OrganizationService) GetRecordByCustomerId(customerId string) (customer EntityCustomer.Customer, err error) {

	const MethodName = "GetRecordByCustomerId"
	var entityBody EntityCustomer.Customer
	entityBody.CustomerId = customerId

	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" TableName::"+entityBody.TableName()+" customerId:::", customerId)
	resp, err := h.Repository.FindOne(entityBody.GetFilterId(), entityBody.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error FindOne:::", err)
		return entityBody, err
	}
	data, err := EntityCustomer.ParseDynamoAtributeToStruct(resp.Item)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:ParseDynamoAtributeToStruct:::", err)
		return data, err
	}
	return data, nil

}

func (h *OrganizationService) GetAllRecords() (entities []EntityCustomer.Customer, err error) {
	const MethodName = "GetAllRecords"
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" id:::", nil)
	entities = []EntityCustomer.Customer{}
	var entity EntityCustomer.Customer

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
			entity, err := EntityCustomer.ParseDynamoAtributeToStruct(value)
			if err != nil {
				logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:ParseDynamoAtributeToStruct::", err)
				return entities, err
			}
			entities = append(entities, entity)
		}
	}
	return entities, err
}

func (h *OrganizationService) CreateRecord(entityBody *EntityCustomer.Customer) (response *events.APIGatewayProxyResponse, err error) {
	const MethodName = "CreateRecord"
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" start:::", entityBody)

	err = h.Rules.Validate(entityBody)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:validate::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusBadRequest, "", err)
	}
	setDefaultValues(entityBody, entityBody.CustomerId)
	entityBody.UpdatedBy = entityBody.CreatedBy
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" entityBody:::", entityBody)
	resp, err := h.Repository.CreateOrUpdate(entityBody.GetMap(), entityBody.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:CreateOrUpdate::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" response:::", resp)

	return ResponseUtil.GetFinalResponse(uuid.Nil, entityBody, http.StatusOK, SUCCESS_CREATED, nil)

}

func (h *OrganizationService) UpdateRecord(entityBody *EntityCustomer.CustomerUpdate) (response *events.APIGatewayProxyResponse, err error) {

	const MethodName = "UpdateRecord"

	var entityCustomer EntityCustomer.Customer

	err = h.Rules.ValidateUpdate(entityBody)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:validate::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, entityBody, http.StatusBadRequest, "", err)
	}
	conditionMap := make(map[string]interface{})
	conditionMap["CustomerId"] = entityBody.CustomerId
	dbData, err := h.Repository.FindOne(conditionMap, entityCustomer.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:FindOne::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, entityBody, http.StatusInternalServerError, "", err)
	}
	found, err := EntityCustomer.ParseDynamoAtributeToStruct(dbData.Item)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:ParseDynamoAtributeToStruct::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, entityBody, http.StatusInternalServerError, "", err)
	}
	setDefaultValues(&found, entityBody.CustomerId)
	found.UpdatedBy = entityBody.UpdatedBy

	found.CustomerName = entityBody.CustomerName
	found.OwnerId = entityBody.OwnerId
	found.Address = entityBody.Address
	found.State = entityBody.State
	found.ZipCode = entityBody.ZipCode
	found.Country = entityBody.Country

	if NETWORK == entityBody.FieldName {
		found.NetWork = entityBody.NetWork
	} else if VENDOR == entityBody.FieldName {
		found.Vendors = entityBody.Vendors
	} else if INTEGRATION == entityBody.FieldName {
		found.Integration = entityBody.Integration
	} else {
		err = errors.New("Invalid request. Fieldname value is incorrect.")
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" Request invalid. Fieldname value is incorrect. ::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, entityBody, http.StatusBadRequest, "", err)
	}

	resp, err := h.Repository.CreateOrUpdate(found.GetMap(), found.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:CreateOrUpdate::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, entityBody, http.StatusInternalServerError, "", err)
	}
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" response:::", resp)

	return ResponseUtil.GetFinalResponse(uuid.Nil, found, http.StatusOK, SUCCESS_UPDATED, nil)

}

func (h *OrganizationService) DeleteRecord(customerId string) (response *events.APIGatewayProxyResponse, err error) {
	const MethodName = "DeleteRecord"

	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" start:::", customerId)
	var entity EntityCustomer.Customer
	entity.CustomerId = customerId

	resp, err := h.Repository.FindOne(entity.GetFilterId(), entity.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:FindOne::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, entity, http.StatusInternalServerError, "", err)
	}
	entity, err = EntityCustomer.ParseDynamoAtributeToStruct(resp.Item)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:FindOne::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, entity, http.StatusInternalServerError, "", err)
	}

	_, err = h.Repository.Delete(entity.GetFilterId(), entity.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:Delete::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, entity, http.StatusInternalServerError, "", err)
	}

	return ResponseUtil.GetFinalResponse(uuid.Nil, entity, http.StatusOK, SUCCESS_DELETED, nil)

}

func setDefaultValues(assesment *EntityCustomer.Customer, ID string) {
	assesment.UpdatedAt = time.Now()
	if ID == "" {
		assesment.CustomerId = uuid.New().String()
		assesment.CreatedAt = time.Now()
	} else {
		assesment.CustomerId = ID
	}

}
