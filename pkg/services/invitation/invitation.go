package invitation

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	//EntityEntities "riscvue.com/pkg/entities"

	EntityInvitation "riscvue.com/pkg/entities/invitation"
	"riscvue.com/pkg/repository/adapter"

	ServiceInterface "riscvue.com/pkg/services"
	"riscvue.com/pkg/utils/logger"
)

const CLASSS_NAME = "InvitationService"

type InvitationService struct {
	ServiceInterface.InvitationInterface
	Repository adapter.Interface
	//Req        *http.Request
}

const INVITEDCUSTOMERID = "InvitedCustomerId"
const INVITEDTO = "InvitedTo"
const SUCCESS_RETRIVED = "Successfully Retrieved."

func NewInvitationService(repository adapter.Interface) ServiceInterface.InvitationInterface {
	return &InvitationService{
		Repository: repository,
		//Req:        req,
	}
}
func (h *InvitationService) FindByCustomerIdAndIsSuperAdmin(customerId string) (entities []EntityInvitation.Invitation, err error) {
	const MethodName = "FindByCustomerIdAndIsSuperAdmin"
	entities = []EntityInvitation.Invitation{}

	var entity EntityInvitation.Invitation
	if customerId != "" {
		filter := expression.Name(INVITEDCUSTOMERID).Equal(expression.Value(customerId))
		condition, err := expression.NewBuilder().WithFilter(filter).Build()
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:NewBuilder:::", err)
			return entities, err

		}

		resp, err := h.Repository.FindAll(condition, entity.TableName())
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:FindAll", err)
			return entities, err

		}

		count := resp.Count
		//check := int64(0)
		if *count == 0 {
			err = errors.New("Item not found")
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:", err)
			return entities, err
		}

		if resp != nil {
			for _, value := range resp.Items {
				entity, err := EntityInvitation.ParseDynamoAtributeToStruct(value)
				if err != nil {
					logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:ParseDynamoAtributeToStruct::", err)
					return entities, err
				}

				logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" entity.InvitedTo::", entity.InvitedTo)
				logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" entity.InvitedCustomerId::", entity.InvitedCustomerId)

				entities = append(entities, entity)
			}
		}
		return entities, nil
	} else {
		err = errors.New("Record not found for given customerId" + customerId)
		return entities, err
	}
}

func (h *InvitationService) IsSuperAdmin(customerId string, userId string) (isSuperAdmin bool, err error) {
	const MethodName = "IsSuperAdmin"
	isSuperAdmin = false

	var entity EntityInvitation.Invitation
	if customerId != "" && userId != "" {
		filter := expression.And(expression.Name(INVITEDCUSTOMERID).Equal(expression.Value(customerId)), expression.Name(INVITEDTO).Equal(expression.Value(userId)))
		condition, err := expression.NewBuilder().WithFilter(filter).Build()
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:NewBuilder:::", err)
			return isSuperAdmin, err

		}

		resp, err := h.Repository.FindAll(condition, entity.TableName())
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:FindAll", err)
			return isSuperAdmin, err

		}

		count := resp.Count
		//check := int64(0)
		if *count == 0 {
			err = errors.New("Item not found")
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:", err)
			return isSuperAdmin, err
		}

		if resp != nil {

			entity, err := EntityInvitation.ParseDynamoAtributeToStruct(resp.Items[0])
			if err != nil {
				logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" error:ParseDynamoAtributeToStruct::", err)
				return isSuperAdmin, err
			}

			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" entity.InvitedTo::", entity.InvitedTo)
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" entity.InvitedCustomerId::", entity.InvitedCustomerId)
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" entity.IsSuperAdmin::", entity.IsSuperAdmin)

			isSuperAdmin = entity.IsSuperAdmin

		}

	}
	return isSuperAdmin, nil
}
