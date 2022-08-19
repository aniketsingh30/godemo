package services

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	EntityAssessment "riscvue.com/pkg/entities/assessment"
	EntityAssessmentSnapshot "riscvue.com/pkg/entities/assessmentsnapshot"
	EntityControl "riscvue.com/pkg/entities/control"
	EntityControlData "riscvue.com/pkg/entities/controldata"
	EntityCustomer "riscvue.com/pkg/entities/customer"
	EntityInvitation "riscvue.com/pkg/entities/invitation"
	EntitySecurityScroreCard "riscvue.com/pkg/entities/securityScoreCard"
	EntityUser "riscvue.com/pkg/entities/user"
)

type AllData struct {
	ControlNumber, ControlFamily, ShortDescription, MappedControl, Requirement string
}

type AssessmentInterface interface {
	Get(id string, history string, createdBy string, name string, customerId string, userId string, all string) (response *events.APIGatewayProxyResponse, err error)
	getRecord(id string) (response *events.APIGatewayProxyResponse, err error)
	getRecordByName(name string) (response *events.APIGatewayProxyResponse, err error)
	findByCreatedBy(createdBy string) (response *events.APIGatewayProxyResponse, err error)
	FindByCustomerId(customerId string) (response *events.APIGatewayProxyResponse, err error)
	FindByUserIdAndCustomerId(userId string, customerId string) (response *events.APIGatewayProxyResponse, err error)
	findAssessmentControls(entities []EntityAssessment.Assessment) (assesmentDatas []EntityAssessment.AssessmentData, err error)
	FindByCreatedByAndCustomerId(createdBy string, customerId string) (response *events.APIGatewayProxyResponse, err error)
	getAllRecords() (response *events.APIGatewayProxyResponse, err error)
	CreateRecord(assessment *EntityAssessment.Assessment) (response *events.APIGatewayProxyResponse, err error)
	UpdateRecord(assessment *EntityAssessment.Assessment) (response *events.APIGatewayProxyResponse, err error)
	DeleteRecord(id string, deletedby string) (response *events.APIGatewayProxyResponse, err error)
	IsAuthorizedUser(customerId string, userId string) (isAuthorized bool, err error)
	UnhandledMethod() (response *events.APIGatewayProxyResponse, err error)
}

// type UserMappingInterface interface {
// 	//Get() (response *events.APIGatewayProxyResponse, err error)
// 	//getRecord(id string) (response *events.APIGatewayProxyResponse, err error)
// 	FindByUserIdAndCustomerId(userid string, customerId string) (response EntityUserMapping.UserMapping, err error)
// 	findByUserId(userid string) (response EntityUserMapping.UserMapping, err error)
// 	//findByCustomerId(customerId string) (response *events.APIGatewayProxyResponse, err error)
// 	//findByUserId(userId string) (response *events.APIGatewayProxyResponse, err error)
// 	//getAllRecords() (response *events.APIGatewayProxyResponse, err error)
// 	CreateRecord() (response *events.APIGatewayProxyResponse, err error)
// 	CreateRecordFromInvitation(userId string, customerId string) (response *events.APIGatewayProxyResponse, err error)
// 	UpdateRecord() (response *events.APIGatewayProxyResponse, err error)
// 	DeleteRecord() (response *events.APIGatewayProxyResponse, err error)
// 	UnhandledMethod() (response *events.APIGatewayProxyResponse, err error)
// 	isRecordExistsByName(name string) (isRecordExists bool, err error)
// }

type AssessmentControlInterface interface {
	FindAssessmentControls(id uuid.UUID, includeHistory bool) (assessmentControls []EntityControl.AssessmentControl, err error)
	CreateControl(assessMentControl EntityControl.AssessmentControl, controlDataItem EntityControlData.ControlData, createdBy string) (entity EntityControl.AssessmentControl, err error)
	DeleteAssessmentControls(assessmentId uuid.UUID) error
}

type AssessmentSnapshotInterface interface {
	FindAssessmentSnapshots(id uuid.UUID) (assessmentSnapshots []EntityAssessmentSnapshot.AssessmentSnapshot, err error)
	DeleteAssessmentSnapshots(assessmentId uuid.UUID) error
}
type ControlDataInterface interface {
	FindControlsByControlDataId(controlDataId string) (entities []EntityControlData.ControlData, err error)
}

type ControlDataInsertInterface interface {
	CreateRecord(controlDataReq *EntityControlData.ControlDataRequest) (controlDatas []EntityControlData.ControlData, err error)
	CreateRecordFromOtherEnv(env string) (controlDatas []EntityControlData.ControlData, err error)
	CreateRecordInternal(excelData AllData) (entity EntityControlData.ControlData, err error)
}

type InvitationInterface interface {
	FindByCustomerIdAndIsSuperAdmin(customerId string) (entities []EntityInvitation.Invitation, err error)
	IsSuperAdmin(customerId string, userId string) (isSuperAdmin bool, err error)
}

type OrganizationInterface interface {
	GetRecordByCustomerId(customerId string) (customer EntityCustomer.Customer, err error)
	GetAllRecords() (entities []EntityCustomer.Customer, err error)
	CreateRecord(mainResponse *EntityCustomer.Customer) (response *events.APIGatewayProxyResponse, err error)
	UpdateRecord(mainResponse *EntityCustomer.CustomerUpdate) (response *events.APIGatewayProxyResponse, err error)
	DeleteRecord(customerId string) (response *events.APIGatewayProxyResponse, err error)
}

type UserInterface interface {
	GetRecordByUserId(userId string) (user EntityUser.User, err error)
	GetAllRecords() (entities []EntityUser.User, err error)
}

type SecurityScoreCardInterface interface {
	Get(domainName string, all string) (response *events.APIGatewayProxyResponse, err error)
	getRecord(id string) (response *events.APIGatewayProxyResponse, err error)
	getRecordByDomain(domainName string) (response *events.APIGatewayProxyResponse, err error)
	getAllRecords() (response *events.APIGatewayProxyResponse, err error)
	CreateRecord(mainResponse *EntitySecurityScroreCard.MainResponse) (response *events.APIGatewayProxyResponse, err error)
	UpdateRecord(mainResponse *EntitySecurityScroreCard.MainResponse) (response *events.APIGatewayProxyResponse, err error)
	DeleteRecord(domainName string, customerId string) (response *events.APIGatewayProxyResponse, err error)
	IsAuthorizedUser(customerId string, userId string) (isAuthorized bool, err error)
	UnhandledMethod() (response *events.APIGatewayProxyResponse, err error)
}

type QualysReportInterface interface {
	Get(domainName string, all string) (response *events.APIGatewayProxyResponse, err error)
	getRecord(id string) (response *events.APIGatewayProxyResponse, err error)
	getRecordByDomain(domainName string) (response *events.APIGatewayProxyResponse, err error)
	getAllRecords() (response *events.APIGatewayProxyResponse, err error)
	CreateRecord(mainResponse *EntitySecurityScroreCard.MainResponse) (response *events.APIGatewayProxyResponse, err error)
	UpdateRecord(mainResponse *EntitySecurityScroreCard.MainResponse) (response *events.APIGatewayProxyResponse, err error)
	DeleteRecord(domainName string, customerId string) (response *events.APIGatewayProxyResponse, err error)
	IsAuthorizedUser(customerId string, userId string) (isAuthorized bool, err error)
	UnhandledMethod() (response *events.APIGatewayProxyResponse, err error)
}
