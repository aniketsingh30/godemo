package assessment

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	EntityAssessment "riscvue.com/pkg/entities/assessment"
	EntityAssessmentRoleMapping "riscvue.com/pkg/entities/assessmentrolemapping"
	EntityAssessmentControl "riscvue.com/pkg/entities/control"
	"riscvue.com/pkg/repository/adapter"
	Rules "riscvue.com/pkg/rules"
	RulesAssessment "riscvue.com/pkg/rules/assessment"
	ServiceInterface "riscvue.com/pkg/services"

	AssessmentControlService "riscvue.com/pkg/services/assessmentcontrol"
	AssessmentSnapshotService "riscvue.com/pkg/services/assessmentsnapshot"
	ControlDataService "riscvue.com/pkg/services/controldata"
	InvitationService "riscvue.com/pkg/services/invitation"
	OrganizationService "riscvue.com/pkg/services/organization"
	UserService "riscvue.com/pkg/services/user"
	"riscvue.com/pkg/utils/logger"
	ResponseUtil "riscvue.com/pkg/utils/response"
)

const CLASSS_NAME = "AssessmentInterface"
const ID = "id"
const RESULT = "result"
const STATUS_CODE = "status"
const SUCCESS_MSG = "successMessage"
const FAILLED_MSG = "failedMessage"
const (
	Accepted    = "Accepted"
	Rejected    = "Rejected"
	CONTRIBUTOR = "Contributor"
	MANAGER     = "Manager"
	ADMIN       = "Admin"
)

var ErrorMethodNotAllowed = "method Not allowed"

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}
type AssesmentService struct {
	ServiceInterface.AssessmentInterface
	Repository adapter.Interface
	Req        *http.Request
	Rules      Rules.Interface
	//UserMappingService        ServiceInterface.UserMappingInterface
	AssessmentControlService  ServiceInterface.AssessmentControlInterface
	ControlDataService        ServiceInterface.ControlDataInterface
	AssessmentSnapshotService ServiceInterface.AssessmentSnapshotInterface
	OrganizationService       ServiceInterface.OrganizationInterface
	UserService               ServiceInterface.UserInterface
	InvitationService         ServiceInterface.InvitationInterface
}

const CREATEDBY = "CreatedBy"
const NAME = "Name"
const INVITEDBY = "InvitedBy"
const CUSTOMERID = "CustomerId"
const USERID = "userId"

const SUCCESS_RETRIVED = "Successfully Retrieved."
const SUCCESS_DELETED = "Successfully Deleted."
const SUCCESS_UPDATED = "Successfully Updated."
const SUCCESS_CREATED = "Successfully Created."

func NewAssessmentService(repository adapter.Interface) ServiceInterface.AssessmentInterface {
	return &AssesmentService{
		Repository: repository,
		//Req:        req,
		Rules: RulesAssessment.NewRules(),
		//UserMappingService:        MappingService.NewUserMappingService(repository),
		AssessmentControlService:  AssessmentControlService.NewAssessmentControlService(repository),
		ControlDataService:        ControlDataService.NewControlDataService(repository),
		AssessmentSnapshotService: AssessmentSnapshotService.NewAssessmentSnapshotService(repository),
		OrganizationService:       OrganizationService.NewOrganizationService(repository),
		UserService:               UserService.NewUserService(repository),
		InvitationService:         InvitationService.NewInvitationService(repository),
	}
}

func (h *AssesmentService) Get(id string, history string, createdBy string, name string, customerId string, userId string, all string) (response *events.APIGatewayProxyResponse, err error) {
	var MethodName = "Get"
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" request:::", createdBy)

	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" request:::", id)
	id = "0e868c26-4f5d-4651-b660-63dcbf73d1e0"
	history = "true"

	if id != "" && history != "" {
		includeHistory, err := strconv.ParseBool(history)

		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:::", err)
			return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusBadRequest, "", err)
		}
		return h.getRecord(id, includeHistory)
	} else if createdBy != "" && customerId != "" {
		return h.FindByCreatedByAndCustomerId(createdBy, customerId)
	} else if userId != "" && customerId != "" {
		return h.FindByUserIdAndCustomerId(userId, customerId)
	} else if createdBy != "" {

		return h.findByCreatedBy(createdBy)
	} else if customerId != "" {

		return h.FindByCustomerId(customerId)
	} else if name != "" {

		return h.getRecordByName(name)
	} else if all != "" {
		return h.getAllRecords()
	} else {
		err := errors.New("Invalid Request.")
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusBadRequest, "", err)
	}
}

func (h *AssesmentService) findAssessmentControls(entities []EntityAssessment.Assessment) (assesmentDatas []EntityAssessment.AssessmentData, err error) {
	const MethodName = "findAssessmentControls"
	assesmentDatas = []EntityAssessment.AssessmentData{}

	for _, assessmentEntity := range entities {
		var assesmentData EntityAssessment.AssessmentData
		assesmentData.Assessment = assessmentEntity
		assessmentControls, err := h.AssessmentControlService.FindAssessmentControls(assessmentEntity.ID, true)
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:::", err)
			return assesmentDatas, err
		}
		assesmentData.AssessmentControls = assessmentControls
		assesmentDatas = append(assesmentDatas, assesmentData)

	}

	return assesmentDatas, nil
}

func (h *AssesmentService) getRecordByName(name string) (response *events.APIGatewayProxyResponse, err error) {

	const MethodName = "getRecordByName"

	isRecordExists, err := h.isRecordExistsByName(name)
	if err != nil {
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}
	// _, err = updateLastActive(h.Repository, h.Req)
	// if err != nil {
	// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:updateLastActive:::", err)
	// 	//return ResponseUtil.GetFinalResponse(ID, data, http.StatusInternalServerError, SUCCESS_RETRIVED, err)
	// }

	return ResponseUtil.GetFinalResponse(uuid.Nil, isRecordExists, http.StatusOK, SUCCESS_RETRIVED, nil)

}

func (h *AssesmentService) isRecordExistsByName(name string) (isRecordExists bool, err error) {

	const MethodName = "isRecordExistsByName"
	var entityBody EntityAssessment.Assessment
	isRecordExists = false

	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" TableName::"+entityBody.TableName()+" name:::", name)
	filter := expression.Name(NAME).Equal(expression.Value(name))
	condition, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:NewBuilder:::", err)
		return isRecordExists, err
	}

	resp, err := h.Repository.FindAll(condition, entityBody.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:FindAll", err)
		return isRecordExists, err
	}
	count := resp.Count

	if *count >= 1 {
		isRecordExists = true

	}
	// _, err = updateLastActive(h.Repository, h.Req)
	// if err != nil {
	// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:updateLastActive:::", err)
	// 	//return ResponseUtil.GetFinalResponse(ID, data, http.StatusInternalServerError, SUCCESS_RETRIVED, err)
	// }

	return isRecordExists, nil

}

func (h *AssesmentService) getRecord(id string, includeHistory bool) (response *events.APIGatewayProxyResponse, err error) {

	const MethodName = "getRecord"
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" id:::", id)
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" includeHistory:::", includeHistory)
	var entityBody EntityAssessment.Assessment
	ID, err := uuid.Parse(id)
	entityBody.ID = ID
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" TableName::"+entityBody.TableName()+" id:::", id)
	resp, err := h.Repository.FindOne(getFilterId(ID), entityBody.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error FindOne:::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}
	data, err := EntityAssessment.ParseDynamoAtributeToStruct(resp.Item)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:ParseDynamoAtributeToStruct:::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}

	var AssessmentData EntityAssessment.AssessmentData
	//entities = append(entities, data)
	controls, err := h.AssessmentControlService.FindAssessmentControls(data.ID, includeHistory)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:findAssessmentControls:::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}
	AssessmentData.Assessment = data
	AssessmentData.AssessmentControls = controls
	// _, err = updateLastActive(h.Repository, h.Req)
	// if err != nil {
	// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:updateLastActive:::", err)
	// 	//return ResponseUtil.GetFinalResponse(ID, data, http.StatusInternalServerError, SUCCESS_RETRIVED, err)
	// }
	return ResponseUtil.GetFinalResponse(uuid.Nil, AssessmentData, http.StatusOK, SUCCESS_RETRIVED, nil)

}

func (h *AssesmentService) findByCreatedBy(createdBy string) (response *events.APIGatewayProxyResponse, err error) {
	const MethodName = "findByCreatedBy"
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" id:::", nil)
	entities := []EntityAssessment.Assessment{}
	var entity EntityAssessment.Assessment
	filter := expression.Name(CREATEDBY).Equal(expression.Value(createdBy))
	condition, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:NewBuilder:::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}

	resp, err := h.Repository.FindAll(condition, entity.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:FindAll", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}

	if resp != nil {
		for _, value := range resp.Items {
			entity, err := EntityAssessment.ParseDynamoAtributeToStruct(value)

			if err != nil {
				logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:ParseDynamoAtributeToStruct::", err)
				return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
			}

			entities = append(entities, entity)
		}
	}
	// assessmentDatas, err := h.findAssessmentControls(entities)
	// if err != nil {
	// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:h.findAssessmentControls:::", err)
	// 	return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	// }
	// _, err = updateLastActive(h.Repository, h.Req)
	// if err != nil {
	// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:updateLastActive:::", err)
	// 	//return ResponseUtil.GetFinalResponse(uuid.Nil, entities, http.StatusInternalServerError, SUCCESS_RETRIVED, err)
	// }
	return ResponseUtil.GetFinalResponse(uuid.Nil, entities, http.StatusOK, SUCCESS_RETRIVED, nil)
}

func (h *AssesmentService) FindByCustomerId(customerId string) (response *events.APIGatewayProxyResponse, err error) {
	const MethodName = "FindByCustomerId"
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" customerId:", customerId)
	entities := []EntityAssessment.Assessment{}
	var entity EntityAssessment.Assessment
	filter := expression.Name(CUSTOMERID).Equal(expression.Value(customerId))

	condition, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:NewBuilder:::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}

	resp, err := h.Repository.FindAll(condition, entity.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:FindAll", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}

	if resp != nil {
		for _, value := range resp.Items {
			entity, err := EntityAssessment.ParseDynamoAtributeToStruct(value)

			if err != nil {
				logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:ParseDynamoAtributeToStruct::", err)
				return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
			}

			entities = append(entities, entity)

		}
	}
	sorted_Assessment := make(EntityAssessment.AssessmentSlice, 0, len(entities))
	for _, assessment := range entities {
		sorted_Assessment = append(sorted_Assessment, assessment)
	}
	sort.Sort(sorted_Assessment)

	// assessmentDatas, err := h.findAssessmentControls(entities)
	// if err != nil {
	// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:h.findAssessmentControls:::", err)
	// 	return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	// }
	// _, err = updateLastActive(h.Repository, h.Req)
	// if err != nil {
	// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:updateLastActive:::", err)
	// 	//return ResponseUtil.GetFinalResponse(uuid.Nil, entities, http.StatusInternalServerError, SUCCESS_RETRIVED, err)
	// }
	return ResponseUtil.GetFinalResponse(uuid.Nil, sorted_Assessment, http.StatusOK, SUCCESS_RETRIVED, nil)
}
func (h *AssesmentService) FindByUserIdAndCustomerId(userId string, customerId string) (response *events.APIGatewayProxyResponse, err error) {
	const MethodName = "FindByUserIdAndCustomerId"
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" userId::customerId:", userId+":"+customerId)
	entities := []EntityAssessment.Assessment{}
	var entity EntityAssessment.Assessment
	filter := expression.Name(CUSTOMERID).Equal(expression.Value(customerId))

	condition, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:NewBuilder:::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}

	resp, err := h.Repository.FindAll(condition, entity.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:FindAll", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}

	if resp != nil {
		for _, value := range resp.Items {
			entity, err := EntityAssessment.ParseDynamoAtributeToStruct(value)

			if err != nil {
				logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:ParseDynamoAtributeToStruct::", err)
				return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
			}
			AssessmentUserMappingRoles := entity.AssessmentUserMappingRoles
			for j := range AssessmentUserMappingRoles {
				if AssessmentUserMappingRoles[j].UserId == userId {
					entities = append(entities, entity)
				}
			}

		}
	}

	// assessmentDatas, err := h.findAssessmentControls(entities)
	// if err != nil {
	// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:h.findAssessmentControls:::", err)
	// 	return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	// }
	// _, err = updateLastActive(h.Repository, h.Req)
	// if err != nil {
	// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:updateLastActive:::", err)
	// 	//return ResponseUtil.GetFinalResponse(uuid.Nil, entities, http.StatusInternalServerError, SUCCESS_RETRIVED, err)
	// }
	return ResponseUtil.GetFinalResponse(uuid.Nil, entities, http.StatusOK, SUCCESS_RETRIVED, nil)
}
func (h *AssesmentService) FindByCreatedByAndCustomerId(createdBy string, customerId string) (response *events.APIGatewayProxyResponse, err error) {
	const MethodName = "FindByCreatedByAndCustomerId"
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" createdBy::customerId:", createdBy+":"+customerId)
	entities := []EntityAssessment.Assessment{}
	var entity EntityAssessment.Assessment
	filter := expression.Name(CREATEDBY).Equal(expression.Value(createdBy))
	filter.And(expression.Name(CUSTOMERID).Equal(expression.Value(customerId)))
	condition, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:NewBuilder:::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}

	resp, err := h.Repository.FindAll(condition, entity.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:FindAll", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}

	if resp != nil {
		for _, value := range resp.Items {
			entity, err := EntityAssessment.ParseDynamoAtributeToStruct(value)

			if err != nil {
				logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:ParseDynamoAtributeToStruct::", err)
				return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
			}

			entities = append(entities, entity)
		}
	}
	// assessmentDatas, err := h.findAssessmentControls(entities)
	// if err != nil {
	// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:h.findAssessmentControls:::", err)
	// 	return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	// }
	// _, err = updateLastActive(h.Repository, h.Req)
	// if err != nil {
	// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:updateLastActive:::", err)
	// 	//return ResponseUtil.GetFinalResponse(uuid.Nil, entities, http.StatusInternalServerError, SUCCESS_RETRIVED, err)
	// }
	return ResponseUtil.GetFinalResponse(uuid.Nil, entities, http.StatusOK, SUCCESS_RETRIVED, nil)
}

func getFilterId(ID uuid.UUID) map[string]interface{} {
	logger.INFO("className="+CLASSS_NAME+" MethodName=getFilterId ID:::", ID.String())
	return map[string]interface{}{"id": ID.String()}
}

func (h *AssesmentService) getAllRecords() (response *events.APIGatewayProxyResponse, err error) {
	const MethodName = "getAllRecords"
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" id:::", nil)
	entities := []EntityAssessment.Assessment{}
	var entity EntityAssessment.Assessment

	filter := expression.Name(CREATEDBY).NotEqual(expression.Value(""))
	condition, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:NewBuilder::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}

	resp, err := h.Repository.FindAll(condition, entity.TableName())
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" response::", response)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName, err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}

	if resp != nil {
		for _, value := range resp.Items {
			entity, err := EntityAssessment.ParseDynamoAtributeToStruct(value)
			if err != nil {
				logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:ParseDynamoAtributeToStruct::", err)
				return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
			}

			entities = append(entities, entity)
		}
	}
	// assessmentDatas, err := h.findAssessmentControls(entities)
	// if err != nil {
	// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:h.findAssessmentControls:::", err)
	// 	return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	// }
	// _, err = updateLastActive(h.Repository, h.Req)
	// if err != nil {
	// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:updateLastActive:::", err)
	// 	//return ResponseUtil.GetFinalResponse(uuid.Nil, entities, http.StatusInternalServerError, SUCCESS_RETRIVED, err)
	// }
	return ResponseUtil.GetFinalResponse(uuid.Nil, entities, http.StatusOK, SUCCESS_RETRIVED, nil)
}

func (h *AssesmentService) CreateRecord(entityBody *EntityAssessment.Assessment) (response *events.APIGatewayProxyResponse, err error) {
	const MethodName = "CreateRecord"
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" start:::", h.Req.Body)
	isAuthorized := false
	// var entityBody EntityAssessment.Assessment
	// if err := json.Unmarshal([]byte(h.Req.Body), &entityBody); err != nil {
	// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:Unmarshal::", err)
	// 	return ResponseUtil.GetFinalResponse(entityBody.ID, nil, http.StatusUnprocessableEntity, "", err)
	// }
	err = h.Rules.Validate(entityBody)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:validate::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusBadRequest, "", err)
	}

	isRecordExists, err := h.isRecordExistsByName(entityBody.Name)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:isRecordExistsByName::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusBadRequest, "", err)
	}
	if isRecordExists {
		err = errors.New("Record with Assessment Name " + entityBody.Name + " already present. Please try with different name.")
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusBadRequest, "", err)
	}
	isAuthorized, err = h.IsAuthorizedUser(entityBody.CustomerId, entityBody.CreatedBy)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:IsAuthorizedUser::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}
	if isAuthorized {
		entityBody.UpdatedBy = entityBody.CreatedBy
		entityBody.Progress = "Todo"
		entityBody.LastActivity = time.Now()
		setDefaultValues(entityBody, uuid.Nil)
		if len(entityBody.AssessmentUserMappingRoles) == 0 {
			entityBody.AssessmentUserMappingRoles = []EntityAssessmentRoleMapping.AssessmentUserMappingRole{}
		}
		// fetching all control data
		entities, err := h.ControlDataService.FindControlsByControlDataId(entityBody.SecurityFrameworkId + "#" + entityBody.MaturityFrameworkId)
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:CreateOrUpdate::", err)
			return ResponseUtil.GetFinalResponse(entityBody.ID, nil, http.StatusInternalServerError, "", err)
		}
		//ends
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" entityBody:::", entityBody)
		resp, err := h.Repository.CreateOrUpdate(entityBody.GetMap(), entityBody.TableName())
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:CreateOrUpdate::", err)
			return ResponseUtil.GetFinalResponse(entityBody.ID, nil, http.StatusInternalServerError, "", err)
		}

		//create controls

		for _, controlDataItem := range entities {

			var assessMentControl EntityAssessmentControl.AssessmentControl

			assessMentControl.AssessmentId = entityBody.ID
			response, err := h.AssessmentControlService.CreateControl(assessMentControl, controlDataItem, entityBody.CreatedBy)
			if err != nil {
				_, _ = h.Repository.Delete(entityBody.GetFilterId(), entityBody.TableName())
				_ = h.AssessmentControlService.DeleteAssessmentControls(entityBody.ID)
				return ResponseUtil.GetFinalResponse(entityBody.ID, entityBody, http.StatusInternalServerError, "", err)
			}
			log.Println(fmt.Printf("response : %v", response))

		}
		// end
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" response:::", resp)
		// _, err = updateLastActive(h.Repository, h.Req)
		// if err != nil {
		// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:updateLastActive:::", err)
		// 	//return ResponseUtil.GetFinalResponse(entityBody.ID, entityBody, http.StatusInternalServerError, SUCCESS_CREATED, err)
		// }
		return ResponseUtil.GetFinalResponse(entityBody.ID, entityBody, http.StatusOK, SUCCESS_CREATED, nil)
	} else {
		err = errors.New("user is not authorized to perform this action. ")
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusUnauthorized, "", err)
	}
}

func (h *AssesmentService) UpdateRecord(entityBody *EntityAssessment.Assessment) (response *events.APIGatewayProxyResponse, err error) {
	const MethodName = "UpdateRecord"
	// var entityBody EntityAssessment.Assessment
	// if err := json.Unmarshal([]byte(h.Req.Body), &entityBody); err != nil {
	// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:Unmarshal::", err)
	// 	return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusUnprocessableEntity, "", err)
	// }
	isAuthorized := false
	//setDefaultValues(&entityBody, entityBody.ID)
	err = h.Rules.ValidateUpdate(entityBody)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:validate::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusBadRequest, "", err)
	}

	isAuthorized, err = h.IsAuthorizedUser(entityBody.CustomerId, entityBody.UpdatedBy)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:IsAuthorizedUser::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}
	conditionMap := make(map[string]interface{})
	conditionMap["id"] = entityBody.ID.String()
	dbData, err := h.Repository.FindOne(conditionMap, entityBody.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:FindOne::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}
	found, err := EntityAssessment.ParseDynamoAtributeToStruct(dbData.Item)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:ParseDynamoAtributeToStruct::", err)
		return ResponseUtil.GetFinalResponse(found.ID, nil, http.StatusInternalServerError, "", err)
	}
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" success:FindOne::", found)
	if found.Name != entityBody.Name {
		isRecordExists, err := h.isRecordExistsByName(entityBody.Name)
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:isRecordExistsByName::", err)
			return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusBadRequest, "", err)
		}
		if isRecordExists {
			err = errors.New("Record with Assessment Name " + entityBody.Name + " already present. Please try with different name.")
			return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusBadRequest, "", err)
		}
	}

	IsAuthorized(found.AssessmentUserMappingRoles, entityBody.UpdatedBy, &isAuthorized)

	if isAuthorized {
		setDefaultValues(&found, found.ID)
		found.UpdatedBy = entityBody.UpdatedBy
		found.Name = entityBody.Name
		//found.MaturityFrameworkId = entityBody.MaturityFrameworkId
		//found.SecurityFrameworkId = entityBody.SecurityFrameworkId
		//AssessmentUserMappingRoles = []EntityAssessmentRoleMapping.AssessmentUserMappingRole{}
		if found.AssessmentUserMappingRoles == nil || len(found.AssessmentUserMappingRoles) == 0 {
			found.AssessmentUserMappingRoles = []EntityAssessmentRoleMapping.AssessmentUserMappingRole{}
		}
		for _, AssesmentUserMappingRoleRequest := range entityBody.AssessmentUserMappingRoles {
			isPresent := false
			for i, AssesmentUserMappingRoleDB := range found.AssessmentUserMappingRoles {
				if AssesmentUserMappingRoleDB.UserId == AssesmentUserMappingRoleRequest.UserId {
					d := &found.AssessmentUserMappingRoles[i]
					d.RoleId = AssesmentUserMappingRoleRequest.RoleId
					isPresent = true
				}

			}
			if !isPresent {
				found.AssessmentUserMappingRoles = append(found.AssessmentUserMappingRoles, AssesmentUserMappingRoleRequest)
			}
		}

		//found.AssessmentUserMappingRoles = entityBody.AssessmentUserMappingRoles
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" found.GetMap():::", found.GetMap())

		resp, err := h.Repository.CreateOrUpdate(found.GetMap(), found.TableName())
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:CreateOrUpdate::", err)
			return ResponseUtil.GetFinalResponse(found.ID, nil, http.StatusInternalServerError, "", err)
		}

		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" response:::", resp)
		// _, err = updateLastActive(h.Repository, h.Req)
		// if err != nil {
		// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:updateLastActive:::", err)
		// 	//return ResponseUtil.GetFinalResponse(found.ID, found, http.StatusInternalServerError, SUCCESS_UPDATED, err)
		// }
		return ResponseUtil.GetFinalResponse(found.ID, found, http.StatusOK, SUCCESS_UPDATED, nil)
	} else {
		err = errors.New("user is not authorized to perform this action. ")
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusUnauthorized, "", err)
	}
}

func (h *AssesmentService) DeleteRecord(id string, deletedBy string) (response *events.APIGatewayProxyResponse, err error) {
	const MethodName = "DeleteRecord"

	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" start:::", id)
	if id == "" || deletedBy == "" {
		err := errors.New("ID or deletedBy  can not be blank.")
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusBadRequest, "", err)
	}
	var entity EntityAssessment.Assessment
	ID, err := uuid.Parse(id)
	entity.ID = ID
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:Parse::", err)
		return ResponseUtil.GetFinalResponse(entity.ID, entity, http.StatusInternalServerError, "", err)
	}

	resp, err := h.Repository.FindOne(entity.GetFilterId(), entity.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:FindOne::", err)
		return ResponseUtil.GetFinalResponse(entity.ID, entity, http.StatusInternalServerError, "", err)
	}
	entity, err = EntityAssessment.ParseDynamoAtributeToStruct(resp.Item)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:FindOne::", err)
		return ResponseUtil.GetFinalResponse(entity.ID, entity, http.StatusInternalServerError, "", err)
	}

	isAuthorized, err := h.IsAuthorizedUser(entity.CustomerId, deletedBy)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:IsAuthorizedUser::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}
	IsAuthorized(entity.AssessmentUserMappingRoles, deletedBy, &isAuthorized)
	//delete controls
	if isAuthorized {
		err = h.AssessmentControlService.DeleteAssessmentControls(entity.ID)
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:DeleteAssessmentControls::", err)
			return ResponseUtil.GetFinalResponse(entity.ID, entity, http.StatusInternalServerError, "", err)
		}
		//delete snapshots
		err = h.AssessmentSnapshotService.DeleteAssessmentSnapshots(entity.ID)
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:DeleteAssessmentSnapshots::", err)
			return ResponseUtil.GetFinalResponse(entity.ID, entity, http.StatusInternalServerError, "", err)
		}
		//delete assessment
		_, err = h.Repository.Delete(entity.GetFilterId(), entity.TableName())
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:Delete::", err)
			return ResponseUtil.GetFinalResponse(entity.ID, entity, http.StatusInternalServerError, "", err)
		}

		// _, err = updateLastActive(h.Repository, h.Req)
		// if err != nil {
		// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:updateLastActive:::", err)
		// 	//return ResponseUtil.GetFinalResponse(entity.ID, entity, http.StatusInternalServerError, SUCCESS_DELETED, err)
		// }
		return ResponseUtil.GetFinalResponse(entity.ID, entity, http.StatusOK, SUCCESS_DELETED, nil)
	} else {
		err = errors.New("user is not authorized to perform this action. ")
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusUnauthorized, "", err)
	}

}

// func updateLastActive(repository adapter.Interface, req events.APIGatewayProxyRequest) (userEntity EntityUser.User, err error) {
// 	const MethodName = "updateLastActive"
// 	var entity EntityUser.User
// 	claims := jwt.MapClaims{}
// 	accessToken := req.Headers["Authorization"]
// 	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
// 		return []byte("<YOUR VERIFICATION KEY>"), nil
// 	})

// 	log.Println(fmt.Printf("Signout  -----err:::%v", err))
// 	log.Println(fmt.Printf("Signout  -----T-O-K-E-N:::%v", token))
// 	userName := claims["email"]
// 	if userName != nil {
// 		username := userName.(string)

// 		conditionMap := make(map[string]interface{})
// 		conditionMap["UserId"] = username
// 		dbData, err := repository.FindOne(conditionMap, entity.TableName())
// 		if err != nil {
// 			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:FindOne:::", err)
// 			return entity, err
// 		}
// 		foundUser, err := EntityUser.ParseDynamoAtributeToStruct(dbData.Item)
// 		if err != nil {
// 			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:ParseDynamoAtributeToStruct:::", err)
// 			return entity, err
// 		}
// 		foundUser.LastActive = time.Now()
// 		_, err = repository.CreateOrUpdate(foundUser.GetMap(), foundUser.TableName())
// 		if err != nil {
// 			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:CreateOrUpate:::", err)
// 			return entity, err
// 		}
// 		return entity, nil
// 	} else {
// 		return entity, errors.New("Invalid  auth token")
// 	}

// }
func (h *AssesmentService) UnhandledMethod() (response *events.APIGatewayProxyResponse, err error) {
	return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusMethodNotAllowed, ErrorMethodNotAllowed, nil)

}

func setDefaultValues(assesment *EntityAssessment.Assessment, ID uuid.UUID) {
	assesment.UpdatedAt = time.Now()
	if ID == uuid.Nil {
		assesment.ID = uuid.New()
		assesment.CreatedAt = time.Now()
	} else {
		assesment.ID = ID
	}
}

func (h *AssesmentService) IsAuthorizedUser(customerId string, userId string) (isAuthorized bool, err error) {

	const MethodName = "IsAuthorizedUser"
	isAuthorized = false
	//validations
	// 1. user super admin or paid
	//get cutomer owner
	//get cutomer paid user or super admin
	// if created by matching with paid, SA or owner user can create Assessment
	customer, err := h.OrganizationService.GetRecordByCustomerId(customerId)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:GetRecordByCustomerId::", err)
		return isAuthorized, err
	}
	if customer.OwnerId == userId {
		isAuthorized = true
	}
	if !isAuthorized {
		isSuper, err := h.InvitationService.IsSuperAdmin(customerId, userId)
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:GetRecordByUserId::", err)
			return isAuthorized, err
		}

		isAuthorized = isSuper

	}
	return isAuthorized, nil
}
func IsAuthorized(assessmentUserMappingRoles []EntityAssessmentRoleMapping.AssessmentUserMappingRole, userId string, isAuthorized *bool) {

	for _, AssesmentUserMappingRoleDB := range assessmentUserMappingRoles {
		if AssesmentUserMappingRoleDB.UserId == userId && (AssesmentUserMappingRoleDB.RoleId == MANAGER || AssesmentUserMappingRoleDB.RoleId == CONTRIBUTOR || AssesmentUserMappingRoleDB.RoleId == ADMIN) {

			*isAuthorized = true
		}

	}
}
