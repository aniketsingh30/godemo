package assessmentsnapshot

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"

	EntityAssessmentSnapshot "riscvue.com/pkg/entities/assessmentsnapshot"
	"riscvue.com/pkg/repository/adapter"
	ServiceInterface "riscvue.com/pkg/services"

	"riscvue.com/pkg/utils/logger"
)

const CLASSS_NAME = "AssessmentSnapshotInterface"

type AssessmentSnapshotService struct {
	ServiceInterface.AssessmentSnapshotInterface
	Repository adapter.Interface
	//Req        *http.Request
}

const ASSESSMENTID = "AssessmentId"

func NewAssessmentSnapshotService(repository adapter.Interface) ServiceInterface.AssessmentSnapshotInterface {
	return &AssessmentSnapshotService{
		Repository: repository,
		//Req:        req,
	}
}

func (h *AssessmentSnapshotService) FindAssessmentSnapshots(assessmentId uuid.UUID) (entities []EntityAssessmentSnapshot.AssessmentSnapshot, err error) {
	entities = []EntityAssessmentSnapshot.AssessmentSnapshot{}
	var entity EntityAssessmentSnapshot.AssessmentSnapshot
	var filter expression.ConditionBuilder
	if assessmentId != uuid.Nil {
		filter = expression.Name(ASSESSMENTID).Equal(expression.Value(assessmentId.String()))

		condition, err := expression.NewBuilder().WithFilter(filter).Build()
		if err != nil {
			return entities, err
		}

		response, err := h.Repository.FindAll(condition, entity.TableName())
		if err != nil {
			return entities, err
		}

		if response != nil {
			for _, value := range response.Items {
				entity, err := EntityAssessmentSnapshot.ParseDynamoAtributeToStruct(value)
				if err != nil {
					return entities, err
				}
				entities = append(entities, entity)
			}
		}
	}
	return entities, nil
}

func (h *AssessmentSnapshotService) DeleteAssessmentSnapshots(assessmentId uuid.UUID) error {
	const MethodName = "DeleteAssessmentSnapshots"
	entities, err := h.FindAssessmentSnapshots(assessmentId)
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:FindAssessmentSnapshots::", err)
		return err
	}

	for _, assesmentSnapshot := range entities {
		_, err := h.Repository.Delete(assesmentSnapshot.GetFilterId(), assesmentSnapshot.TableName())
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:Delete::", err)
			return err
		}

	}
	return nil
}
