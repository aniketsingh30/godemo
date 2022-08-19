package securityScoreCard

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	EntitySecurityScroreCard "riscvue.com/pkg/entities/securityScoreCard"
	EntityUser "riscvue.com/pkg/entities/user"
	"riscvue.com/pkg/repository/adapter"
	Rules "riscvue.com/pkg/rules"
	RulesSecurityScoreCard "riscvue.com/pkg/rules/securityScoreCard"
	ServiceInterface "riscvue.com/pkg/services"
	InvitationService "riscvue.com/pkg/services/invitation"
	OrganizationService "riscvue.com/pkg/services/organization"
	UserService "riscvue.com/pkg/services/user"
	"riscvue.com/pkg/utils/logger"
	ResponseUtil "riscvue.com/pkg/utils/response"
)

const CLASSS_NAME = "SecurityScoreCardService"
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
type SecurityScoreCardService struct {
	ServiceInterface.SecurityScoreCardInterface
	Repository adapter.Interface
	Req        *http.Request
	Rules      Rules.Interface

	OrganizationService ServiceInterface.OrganizationInterface
	UserService         ServiceInterface.UserInterface
	InvitationService   ServiceInterface.InvitationInterface
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

func NewSecurityScoreCardService(repository adapter.Interface) ServiceInterface.SecurityScoreCardInterface {
	return &SecurityScoreCardService{
		Repository: repository,
		//Req:        req,
		Rules: RulesSecurityScoreCard.NewRules(),

		OrganizationService: OrganizationService.NewOrganizationService(repository),
		UserService:         UserService.NewUserService(repository),
		InvitationService:   InvitationService.NewInvitationService(repository),
	}
}

func (h *SecurityScoreCardService) Get(domainName string, all string) (response *events.APIGatewayProxyResponse, err error) {
	var MethodName = "Get"
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" request:::", h.Req)

	if domainName != "" {

		return h.getRecordByDomainName(domainName)
	} else if all != "" {
		return h.getAllRecords()
	} else {
		err := errors.New("Invalid Request.")
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusBadRequest, "", err)
	}
}

func (h *SecurityScoreCardService) getRecordByDomainName(domainName string) (response *events.APIGatewayProxyResponse, err error) {

	const MethodName = "getRecordByDomainName"
	var entityBody EntitySecurityScroreCard.MainResponse
	entityBody.DomainName = domainName
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" TableName::"+entityBody.TableName()+" domainName:::", domainName)

	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}
	conditionMap := make(map[string]interface{})
	conditionMap["DomainName"] = domainName
	resp, err := h.Repository.FindOne(conditionMap, entityBody.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error FindOne:::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}
	data, err := EntitySecurityScroreCard.ParseDynamoAtributeToStruct(resp.Item)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:ParseDynamoAtributeToStruct:::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}

	// //_, err = updateLastActive(h.Repository, h.Req)
	// if err != nil {
	// 	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:updateLastActive:::", err)

	// }
	return ResponseUtil.GetFinalResponse(uuid.Nil, data, http.StatusOK, SUCCESS_RETRIVED, nil)
}

func (h *SecurityScoreCardService) getAllRecords() (response *events.APIGatewayProxyResponse, err error) {
	const MethodName = "getAllRecords"
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" id:::", nil)
	entities := []EntitySecurityScroreCard.MainResponse{}
	var entity EntitySecurityScroreCard.MainResponse

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
			entity, err := EntitySecurityScroreCard.ParseDynamoAtributeToStruct(value)
			if err != nil {
				logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:ParseDynamoAtributeToStruct::", err)
				return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
			}

			entities = append(entities, entity)
		}
	}

	return ResponseUtil.GetFinalResponse(uuid.Nil, entities, http.StatusOK, SUCCESS_RETRIVED, nil)
}

func (h *SecurityScoreCardService) CreateRecord(entityBody *EntitySecurityScroreCard.MainResponse) (response *events.APIGatewayProxyResponse, err error) {
	const MethodName = "CreateRecord"
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" start:::", h.Req.Body)
	isAuthorized := false

	err = h.Rules.Validate(entityBody)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:validate::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusBadRequest, "", err)
	}
	url := os.Getenv("EXTERNAL_SERVICE_URI")
	if url == "" {
		url = "https://platform-api.securityscorecard.io/companies/makenacap.com"
	} else {
		url = url + "/makenacap.com"
	}
	cleint := &http.Client{}
	byteArry, err := GetData(url, *cleint)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error ScoreResponse:::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}
	var ScoreResponse EntitySecurityScroreCard.ScoreResponse
	json.Unmarshal(byteArry, &ScoreResponse)
	//fmt.Printf("API Response as struct %+v\n", ScoreResponse)
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:ScoreResponse::", ScoreResponse)
	url = url + "/factors"
	byteArry1, err := GetData(url, *cleint)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error Factors:::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}
	var Factors EntitySecurityScroreCard.Factors
	json.Unmarshal(byteArry1, &Factors)
	//fmt.Printf("API Response as struct %+v\n", ScoreResponse)
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:Factors::", Factors)

	entityBody.ScoreResponse = ScoreResponse
	entityBody.Factors = Factors
	isAuthorized, err = h.IsAuthorizedUser(entityBody.CustomerId, entityBody.CreatedBy)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:IsAuthorizedUser::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}
	if isAuthorized {
		entityBody.UpdatedBy = entityBody.CreatedBy
		setDefaultValues(entityBody, uuid.Nil)
		resp, err := h.Repository.CreateOrUpdate(entityBody.GetMap(), entityBody.TableName())
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:CreateOrUpdate::", err)
			return ResponseUtil.GetFinalResponse(entityBody.ID, nil, http.StatusInternalServerError, "", err)
		}

		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" response:::", resp)

		return ResponseUtil.GetFinalResponse(entityBody.ID, entityBody, http.StatusOK, SUCCESS_CREATED, nil)
	} else {
		err = errors.New("user is not authorized to perform this action. ")
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusUnauthorized, "", err)
	}
}

func (h *SecurityScoreCardService) UpdateRecord(entityBody *EntitySecurityScroreCard.MainResponse) (response *events.APIGatewayProxyResponse, err error) {
	const MethodName = "UpdateRecord"

	isAuthorized := false

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
	conditionMap["DomainName"] = entityBody.DomainName
	dbData, err := h.Repository.FindOne(conditionMap, entityBody.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:FindOne::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}
	found, err := EntitySecurityScroreCard.ParseDynamoAtributeToStruct(dbData.Item)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:ParseDynamoAtributeToStruct::", err)
		return ResponseUtil.GetFinalResponse(found.ID, nil, http.StatusInternalServerError, "", err)
	}
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" success:FindOne::", found)

	if isAuthorized {
		setDefaultValues(&found, found.ID)
		found.UpdatedBy = entityBody.UpdatedBy

		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" found.GetMap():::", found.GetMap())

		resp, err := h.Repository.CreateOrUpdate(found.GetMap(), found.TableName())
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:CreateOrUpdate::", err)
			return ResponseUtil.GetFinalResponse(found.ID, nil, http.StatusInternalServerError, "", err)
		}

		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" response:::", resp)

		return ResponseUtil.GetFinalResponse(found.ID, found, http.StatusOK, SUCCESS_UPDATED, nil)
	} else {
		err = errors.New("user is not authorized to perform this action. ")
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusUnauthorized, "", err)
	}
}

func (h *SecurityScoreCardService) DeleteRecord(domainName string, deletedBy string) (response *events.APIGatewayProxyResponse, err error) {
	const MethodName = "DeleteRecord"

	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" start:::", domainName)
	if domainName == "" || deletedBy == "" {
		err := errors.New("domainName or deletedBy  can not be blank.")
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusBadRequest, "", err)
	}
	var entity EntitySecurityScroreCard.MainResponse
	entity.DomainName = domainName
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:Parse::", err)
		return ResponseUtil.GetFinalResponse(entity.ID, entity, http.StatusInternalServerError, "", err)
	}

	resp, err := h.Repository.FindOne(entity.GetFilterId(), entity.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:FindOne::", err)
		return ResponseUtil.GetFinalResponse(entity.ID, entity, http.StatusInternalServerError, "", err)
	}
	entity, err = EntitySecurityScroreCard.ParseDynamoAtributeToStruct(resp.Item)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:FindOne::", err)
		return ResponseUtil.GetFinalResponse(entity.ID, entity, http.StatusInternalServerError, "", err)
	}

	isAuthorized, err := h.IsAuthorizedUser(entity.CustomerId, deletedBy)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:IsAuthorizedUser::", err)
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusInternalServerError, "", err)
	}

	//delete controls
	if isAuthorized {

		_, err = h.Repository.Delete(entity.GetFilterId(), entity.TableName())
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:Delete::", err)
			return ResponseUtil.GetFinalResponse(entity.ID, entity, http.StatusInternalServerError, "", err)
		}

		return ResponseUtil.GetFinalResponse(entity.ID, entity, http.StatusOK, SUCCESS_DELETED, nil)
	} else {
		err = errors.New("user is not authorized to perform this action. ")
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusUnauthorized, "", err)
	}

}

func updateLastActive(repository adapter.Interface, req events.APIGatewayProxyRequest) (userEntity EntityUser.User, err error) {
	const MethodName = "updateLastActive"
	var entity EntityUser.User
	claims := jwt.MapClaims{}
	accessToken := req.Headers["Authorization"]
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("<YOUR VERIFICATION KEY>"), nil
	})

	log.Println(fmt.Printf("Signout  -----err:::%v", err))
	log.Println(fmt.Printf("Signout  -----T-O-K-E-N:::%v", token))
	userName := claims["email"]
	if userName != nil {
		username := userName.(string)

		conditionMap := make(map[string]interface{})
		conditionMap["UserId"] = username
		dbData, err := repository.FindOne(conditionMap, entity.TableName())
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:FindOne:::", err)
			return entity, err
		}
		foundUser, err := EntityUser.ParseDynamoAtributeToStruct(dbData.Item)
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:ParseDynamoAtributeToStruct:::", err)
			return entity, err
		}
		foundUser.LastActive = time.Now()
		_, err = repository.CreateOrUpdate(foundUser.GetMap(), foundUser.TableName())
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:CreateOrUpate:::", err)
			return entity, err
		}
		return entity, nil
	} else {
		return entity, errors.New("Invalid  auth token")
	}

}
func (h *SecurityScoreCardService) UnhandledMethod() (response *events.APIGatewayProxyResponse, err error) {
	return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusMethodNotAllowed, ErrorMethodNotAllowed, nil)

}

func setDefaultValues(securityCard *EntitySecurityScroreCard.MainResponse, ID uuid.UUID) {
	securityCard.UpdatedAt = time.Now()
	if ID == uuid.Nil {
		securityCard.ID = uuid.New()
		securityCard.CreatedAt = time.Now()
	} else {
		securityCard.ID = ID
	}
}

func (h *SecurityScoreCardService) IsAuthorizedUser(customerId string, userId string) (isAuthorized bool, err error) {

	const MethodName = "IsAuthorizedUser"
	isAuthorized = false
	//validations
	// 1. user super admin or paid
	//get cutomer owner
	//get cutomer paid user or super admin
	// if created by matching with paid, SA or owner user can create SecurityScore
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

func GetData(url string, client http.Client) ([]byte, error) {

	//client1 := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Token ZJfQ888QHGnCkZySsLqkOp1DSklH")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)

}
