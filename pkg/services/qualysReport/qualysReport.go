package qualysReport

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	xj "github.com/basgys/goxml2json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	EntityQualysReport "riscvue.com/pkg/entities/qualysReport"
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

const CLASSS_NAME = "QualysReportService"
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
type QualysReportService struct {
	ServiceInterface.QualysReportInterface
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

func NewQualysReportService(repository adapter.Interface) ServiceInterface.QualysReportInterface {
	return &QualysReportService{
		Repository: repository,
		//Req:        req,
		Rules: RulesSecurityScoreCard.NewRules(),

		OrganizationService: OrganizationService.NewOrganizationService(repository),
		UserService:         UserService.NewUserService(repository),
		InvitationService:   InvitationService.NewInvitationService(repository),
	}
}

func (h *QualysReportService) Get(domainName string, all string) (response *events.APIGatewayProxyResponse, err error) {
	var MethodName = "Get"
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" request:::", h.Req)
	get()
	if domainName != "" {

		return h.getRecordByDomainName(domainName)
	} else if all != "" {
		return h.getAllRecords()
	} else {
		err := errors.New("Invalid Request.")
		return ResponseUtil.GetFinalResponse(uuid.Nil, nil, http.StatusBadRequest, "", err)
	}
}

func (h *QualysReportService) getRecordByDomainName(domainName string) (response *events.APIGatewayProxyResponse, err error) {

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

func (h *QualysReportService) getAllRecords() (response *events.APIGatewayProxyResponse, err error) {
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

func (h *QualysReportService) CreateRecord(entityBody *EntitySecurityScroreCard.MainResponse) (response *events.APIGatewayProxyResponse, err error) {
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

func (h *QualysReportService) UpdateRecord(entityBody *EntitySecurityScroreCard.MainResponse) (response *events.APIGatewayProxyResponse, err error) {
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

func (h *QualysReportService) DeleteRecord(domainName string, deletedBy string) (response *events.APIGatewayProxyResponse, err error) {
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
func (h *QualysReportService) UnhandledMethod() (response *events.APIGatewayProxyResponse, err error) {
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

func (h *QualysReportService) IsAuthorizedUser(customerId string, userId string) (isAuthorized bool, err error) {

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

	get()
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

var qauth64 string
var qApi *string
var ips *string
var assets *string
var hosts *string
var authFile *string
var authCli *string

func get() {
	//
	// QualysGuard API example: pull Vulnerability details from hosts scanned with Qualys
	//

	// Potentially useful guides to the Qualys API
	//
	// https://www.qualys.com/docs/qualys-api-vmpc-user-guide.pdf
	// https://community.qualys.com/docs/DOC-4523-qualys-api-client-examples#jive_content_id_Go_Language_Example
	// https://www.qualys.com/docs/qualys-asset-management-tagging-api-v2-user-guide.pdf
	// https://community.qualys.com/thread/18542-get-list-of-all-applications-on-authenticated-hosts-via-api
	//

	//
	// The Qualys structs were generated with Zek (https://github.com/miku/zek) using raw XML captured from exploring
	// the Qualys API manually:
	//
	//	`cat raw_output.xml | ./zek > golangstructname.go`
	//

	// Parse the CLI arguments, check for validity

	qApi = flag.String("Api", "https://qualysapi.qg2.apps.qualys.com:443", "Point to the Qualys API endpoint (Review Qualys documentation to see which endpoint you should use)")
	ips = flag.String("HostIps", "89.99.99.101,53.54.55.55", "Enter a comma-delimited list of host IP addresses that have been scanned by Qualys")
	assets = flag.String("AssetIds", "", "Enter a comma-delimited list of Qualys Asset IDs")
	hosts = flag.String("HostIds", "", "HostIDs passed as a comma-delimited list of Qualys host IDs")
	authFile = flag.String("AuthFile", "credentials", "Auth file that contains username/password authentication. Should be in the format of `USERNAME:PASSWORD`\nTo generate a mock/example file, pass in a value of GENERATE_SAMPLE like so:\n\t-AuthFile=GENERATE_SAMPLE\n\t(a file named `credentials` will be created")
	authCli = flag.String("", "", "Can be used in place of the 'AuthFile' if you'd rather specify the username and password at the CLI.\nEx:\n\t-AuthCLI=USERNAME:PASSWORD\n\n\tQualys uses HTTP Basic authentication, which is why we take the credentials in this format")
	flag.Parse()

	if *authFile != "" && *authCli != "" {
		fmt.Println("Specify either AuthFile -or- AuthCLI. Both cannot be specified at once")
		os.Exit(127)
	}

	if (*ips != "" && *assets != "") || (*ips != "" && *hosts != "") || (*hosts != "" && *assets != "") {
		fmt.Println("Specify only ONE of the following CLI options:\n\tHostIps\n\tHostIds\n\tAssetIds")
		os.Exit(126)
	}

	if *authFile == "GENERATE_SAMPLE" {
		ioutil.WriteFile("credentials", []byte("cffer2au:mab&RAs49p"), 0600)
		os.Exit(0)
	}

	// Parse authentication
	var authBytes []byte
	if *authFile != "" {
		authFileBytes, authFileBytesErr := ioutil.ReadFile(*authFile)
		if authFileBytesErr != nil {
			panic(authFileBytesErr)
		}

		authBytes = authFileBytes
	} else {
		authBytes = []byte(*authCli)
	}
	qauth64 = base64.StdEncoding.EncodeToString(authBytes)

	//add ip adress
	if *ips != "" {
		// addIpsOutput := addIPAddress(*ips)
		// log.Println(addIpsOutput)
		//  convertXMlToJson(addIpsOutput)
		// createAsses := createAssest(*ips)
		// log.Println(createAsses)
		// byteArray, err := convertXMlToJson(createAsses)
		// if err != nil {
		// 	logger.INFO("className="+CLASSS_NAME+" MethodName=convertXMlToJson err:Unmarshal::", err)
		// }
		// var entityBody EntityQualysReport.QualysResponse
		// // data, err := json.Marshal(byteArray.String())
		// // if err != nil {
		// // 	logger.INFO("className="+CLASSS_NAME+" MethodName=convertXMlToJson err:Unmarshal::", err)
		// // }
		// if err := json.Unmarshal(byteArray.Bytes(), &entityBody); err != nil {
		// 	logger.INFO("className="+CLASSS_NAME+" MethodName=convertXMlToJson err:Unmarshal::", err)

		// }
		// logger.INFO("className="+CLASSS_NAME+" MethodName=convertXMlToJson success:Unmarshal::", entityBody)
		launchScan := launchScan(qauth64)

		log.Println(launchScan)
		byteArray, err := convertXMlToJson(launchScan)
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName=convertXMlToJson err:Unmarshal::", err)
		}
		var entityBody EntityQualysReport.QualysResponse
		// data, err := json.Marshal(byteArray.String())
		// if err != nil {
		// 	logger.INFO("className="+CLASSS_NAME+" MethodName=convertXMlToJson err:Unmarshal::", err)
		// }
		if err := json.Unmarshal(byteArray.Bytes(), &entityBody); err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName=convertXMlToJson err:Unmarshal::", err)

		}
		logger.INFO("className="+CLASSS_NAME+" MethodName=convertXMlToJson success:Unmarshal::", entityBody)
		// editAssest := editAssest(*ips)
		// log.Println(editAssest)
		// convertXMlToJson(editAssest)
	}
	//
	// Get the list of Asset IDs
	//var assetIdentifiers []string
	// if *ips != "" {
	// 	// Get the list of hosts by IP address
	// 	hostsOutput := searchHostsByIps(*ips)
	// 	var hostIds string
	// 	for _, host := range hostsOutput.RESPONSE.HOSTLIST.HOST {
	// 		hostIds += "," + host.ID.Text
	// 	}
	// 	hostIds = strings.TrimPrefix(hostIds, ",")

	// 	// Convert the hosts into AssetIDs
	// 	serviceResponse := searchHostsQps(hostIds)
	// 	for _, hostAsset := range serviceResponse.Data.HostAsset {
	// 		assetIdentifiers = append(assetIdentifiers, hostAsset.ID.Text)
	// 	}

	// } else if *hosts != "" {
	// 	// Convert the hostIDs into AssetIDs
	// 	serviceResponse := searchHostsQps(*hosts)
	// 	for _, hostAsset := range serviceResponse.Data.HostAsset {
	// 		assetIdentifiers = append(assetIdentifiers, hostAsset.ID.Text)
	// 	}

	// } else if *assets != "" {
	// 	// Already have the Asset identifiers, just need to get them into a slice
	// 	assetIdentifiers = strings.Split(*assets, ",")
	// }

	//
	// Get detailed host information using the asset identifiers
	//		Build the list of QIDs
	// var qids string
	// for _, ai := range assetIdentifiers {
	// 	hostDetail := hostAssetDetails(ai)
	// 	for _, v := range hostDetail.Data.HostAsset.Vuln.List.HostAssetVuln {
	// 		qids += v.Qid.Text + ","
	// 		fmt.Println(v)
	// 	}
	// }
	// qids = strings.TrimSuffix(qids, ",")

	// //
	// // Get the vulnerability details based on the list of QIDs
	// //		Print the details
	// details := vulnerabilityDetails(qids)
	// for _, vDetails := range details.RESPONSE.VULNLIST.VULN {
	// 	fmt.Println(vDetails)
	// }

	// // Print the raw list of qids
	// fmt.Println("===")
	// fmt.Println(qids)
}

// //
// // Take a comma-delimited list of Qualys QIDs and get the vulnerability details
// func vulnerabilityDetails(qids string) qualys.KNOWLEDGEBASEVULNLISTOUTPUT {
// 	vulnDetailsRaw := qApiCallXml("GET", qauth64, "/api/2.0/fo/knowledge_base/vuln/?action=list&details=All&ids="+qids)

// 	var vulnlistdetails qualys.KNOWLEDGEBASEVULNLISTOUTPUT
// 	xml.Unmarshal(vulnDetailsRaw, &vulnlistdetails)

// 	return vulnlistdetails
// }

// //
// // Returns detailed host information for a specific host
// func hostAssetDetails(assetIdentifier string) qualys.ServiceResponseSingleAsset {
// 	hostAssetDetailsRaw := qApiCallXml("GET", qauth64, "/qps/rest/2.0/get/am/hostasset/"+assetIdentifier)

// 	var hostAssetDetailsServiceResponse qualys.ServiceResponseSingleAsset
// 	xml.Unmarshal(hostAssetDetailsRaw, &hostAssetDetailsServiceResponse)

// 	return hostAssetDetailsServiceResponse
// }

//add Ips
func addIPAddress(ipAddresses string) *strings.Reader {
	// Remove any potential whitespace
	strings.Replace(ipAddresses, " ", "", -1)
	// Ensure there are no leading or trailing commas
	strings.Trim(ipAddresses, ",")

	// add IPS in  Qualys
	hostsRawBytes := qApiCallXml("POST", qauth64, "/api/2.0/fo/asset/ip/?action=add", ipAddresses)
	xml := strings.NewReader(string(hostsRawBytes))
	// var hostsListOutput qualys.ADDIPOUTPUTJSON
	// json.Unmarshal(hostsRawBytes, &hostsListOutput)

	// return hostsListOutput
	return xml
}

//add Ips
func createAssest(ipAddresses string) *strings.Reader {
	// Remove any potential whitespace
	strings.Replace(ipAddresses, " ", "", -1)
	// Ensure there are no leading or trailing commas
	strings.Trim(ipAddresses, ",")
	data := url.Values{}
	data.Set("title", "Test-Group10")
	data.Set("ips", ipAddresses)

	// add IPS in  Qualys
	//hostsRawBytes := qApiCallXml("POST", qauth64, "", data)
	qApiUrl := *qApi + "/api/2.0/fo/asset/group/?action=add"
	qHttpReq, qHttpReqErr := http.NewRequest("POST", qApiUrl, strings.NewReader(data.Encode()))
	qHttpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	qHttpReq.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	if qHttpReqErr != nil {
		panic(qHttpReqErr)
	}

	// Populate the headers
	qHttpReq.Header.Add("X-requested-with", "qualys-integrator")
	qHttpReq.Header.Add("Authorization", "Basic "+qauth64)

	qHttpResp, qHttpRespErr := http.DefaultClient.Do(qHttpReq)
	if qHttpRespErr != nil {
		panic(qHttpRespErr)
	}
	qHttpRespBytes, qHttpRespStrErr := ioutil.ReadAll(qHttpResp.Body)
	if qHttpRespStrErr != nil {
		panic(qHttpRespStrErr)
	}

	xml := strings.NewReader(string(qHttpRespBytes))
	// var hostsListOutput qualys.ADDIPOUTPUTJSON
	// json.Unmarshal(hostsRawBytes, &hostsListOutput)

	// return hostsListOutput
	return xml
}

//edit Assest
func editAssest(ipAddresses string) *strings.Reader {
	// Remove any potential whitespace
	strings.Replace(ipAddresses, " ", "", -1)
	// Ensure there are no leading or trailing commas
	strings.Trim(ipAddresses, ",")
	data := url.Values{}
	data.Set("title", "Test-Group4")
	data.Set("ips", ipAddresses)
	data.Set("id", "1234")

	// add IPS in  Qualys
	//hostsRawBytes := qApiCallXml("POST", qauth64, "", data)
	qApiUrl := *qApi + "/api/2.0/fo/asset/group/?action=edit"
	qHttpReq, qHttpReqErr := http.NewRequest("POST", qApiUrl, strings.NewReader(data.Encode()))
	qHttpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	qHttpReq.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	if qHttpReqErr != nil {
		panic(qHttpReqErr)
	}

	// Populate the headers
	qHttpReq.Header.Add("X-requested-with", "qualys-integrator")
	qHttpReq.Header.Add("Authorization", "Basic "+qauth64)

	qHttpResp, qHttpRespErr := http.DefaultClient.Do(qHttpReq)
	if qHttpRespErr != nil {
		panic(qHttpRespErr)
	}
	qHttpRespBytes, qHttpRespStrErr := ioutil.ReadAll(qHttpResp.Body)
	if qHttpRespStrErr != nil {
		panic(qHttpRespStrErr)
	}

	xml := strings.NewReader(string(qHttpRespBytes))
	// var hostsListOutput qualys.ADDIPOUTPUTJSON
	// json.Unmarshal(hostsRawBytes, &hostsListOutput)

	// return hostsListOutput
	return xml
}

func launchScan(qauth64 string) *strings.Reader {

	data := url.Values{}

	data.Set("ip", "89.99.99.101")

	data.Set("scan_title", "My+Test+Vul+Scan")
	data.Set("client_name", "PBS")
	data.Set("option_id", "314051")

	qApiUrl := *qApi + "/api/2.0/fo/scan/?action=launch"

	qHttpReq, qHttpReqErr := http.NewRequest("POST", qApiUrl, strings.NewReader(data.Encode()))
	qHttpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	qHttpReq.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	if qHttpReqErr != nil {
		panic(qHttpReqErr)
	}

	// Populate the headers
	qHttpReq.Header.Add("X-requested-with", "qualys-integrator")
	qHttpReq.Header.Add("Authorization", "Basic "+qauth64)

	qHttpResp, qHttpRespErr := http.DefaultClient.Do(qHttpReq)
	if qHttpRespErr != nil {
		panic(qHttpRespErr)
	}
	qHttpRespBytes, qHttpRespStrErr := ioutil.ReadAll(qHttpResp.Body)
	if qHttpRespStrErr != nil {
		panic(qHttpRespStrErr)
	}

	xml := strings.NewReader(string(qHttpRespBytes))
	// var hostsListOutput qualys.ADDIPOUTPUTJSON
	// json.Unmarshal(hostsRawBytes, &hostsListOutput)

	// return hostsListOutput
	return xml
}

// //
// // Get the HostListOutput from Qualys by searching via IP Addresses
// func searchHostsByIps(ipAddresses string) qualys.HOSTLISTOUTPUT {
// 	// Remove any potential whitespace
// 	strings.Replace(ipAddresses, " ", "", -1)
// 	// Ensure there are no leading or trailing commas
// 	strings.Trim(ipAddresses, ",")

// 	// Get the raw host data back from Qualys
// 	hostsRawBytes := qApiCallXml("GET", qauth64, "/api/2.0/fo/asset/host/?action=list&ips="+ipAddresses)
// 	var hostsListOutput qualys.HOSTLISTOUTPUT
// 	xml.Unmarshal(hostsRawBytes, &hostsListOutput)

// 	return hostsListOutput
// }

// //
// // Get a ServiceResponse which includes host details
// func searchHostsQps(hostIds string) qualys.ServiceResponse {
// 	searchResults := qApiCallXml("POST", qauth64, "/qps/rest/2.0/search/am/hostasset", `<ServiceRequest>
//   <filters>
//     <Criteria field="qwebHostId" operator="IN">`+hostIds+`</Criteria>
//   </filters>
// </ServiceRequest>`)

// 	var response qualys.ServiceResponse
// 	xml.Unmarshal(searchResults, &response)

// 	return response
// }

// //
// Bare-bones function that calls the specified API endpoint and returns the results as a raw byte array
func qApiCallXml(method string, auth64 string, url string, postbody ...string) []byte {
	// Get an HTTP Request
	qApiUrl := *qApi + url
	var qHttpReq *http.Request
	var qHttpReqErr error
	if len(postbody) > 0 {
		// Only add a postbody if one is supplied
		qHttpReq, qHttpReqErr = http.NewRequest(method, qApiUrl, bytes.NewBuffer([]byte(postbody[0])))
	} else {
		qHttpReq, qHttpReqErr = http.NewRequest(method, qApiUrl, nil)
	}
	if qHttpReqErr != nil {
		panic(qHttpReqErr)
	}

	// Populate the headers
	qHttpReq.Header.Add("X-requested-with", "qualys-integrator")
	qHttpReq.Header.Add("Authorization", "Basic "+auth64)

	qHttpResp, qHttpRespErr := http.DefaultClient.Do(qHttpReq)
	if qHttpRespErr != nil {
		panic(qHttpRespErr)
	}
	qHttpRespBytes, qHttpRespStrErr := ioutil.ReadAll(qHttpResp.Body)
	if qHttpRespStrErr != nil {
		panic(qHttpRespStrErr)
	}

	return qHttpRespBytes
}

func convertXMlToJson(xml *strings.Reader) (*bytes.Buffer, error) {
	// Extract data from restful.Request

	// Extract data from restful.Request

	// Convert
	json, err := xj.Convert(xml)
	if err != nil {
		// Oops...
	}

	// ... Use JSON ...
	log.Println(json)
	return json, err
}
