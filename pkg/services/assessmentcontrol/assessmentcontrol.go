package assessmentcontrol

import (
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"

	Entities "riscvue.com/pkg/entities"
	EntityControl "riscvue.com/pkg/entities/control"
	EntityControlData "riscvue.com/pkg/entities/controldata"
	"riscvue.com/pkg/repository/adapter"
	ServiceInterface "riscvue.com/pkg/services"
	"riscvue.com/pkg/utils/logger"
)

const CLASSS_NAME = "AssesmentControlService"

type AssesmentControlService struct {
	ServiceInterface.AssessmentControlInterface
	Repository adapter.Interface
	//	Req        *http.Request
}

const ASSESSMENTID = "AssessmentId"

func NewAssessmentControlService(repository adapter.Interface) ServiceInterface.AssessmentControlInterface {
	return &AssesmentControlService{
		Repository: repository,
		//Req:        req,
	}
}

func (h *AssesmentControlService) FindAssessmentControls(assessmentId uuid.UUID, includeHistory bool) (entities []EntityControl.AssessmentControl, err error) {
	const MethodName = "FindAssessmentControls"
	entities = []EntityControl.AssessmentControl{}
	var entity EntityControl.AssessmentControl
	//var filter expression.ConditionBuilder
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+"assessmentId::", assessmentId)
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+"includeHistory::", includeHistory)
	if assessmentId != uuid.Nil {
		//filter = expression.Name(ASSESSMENTID).Equal(expression.Value(assessmentId.String()))

		proj := expression.NamesList(expression.Name("id"),
			expression.Name("createdAt"), expression.Name("updatedAt"),
			expression.Name("createdBy"), expression.Name("updatedBy"),
			expression.Name("AssessmentId"), expression.Name("ControlNumber"),
			expression.Name("ControlFamily"), expression.Name("ShortDescription"),
			expression.Name("MappedControl"), expression.Name("Requirement"),
			expression.Name("IsEnable"), expression.Name("Status"),
			expression.Name("CurrentScore"), expression.Name("TargetScore"),
			expression.Name("DueDate"), expression.Name("AssignedTo"),
			expression.Name("Likelihood"), expression.Name("Impact"),
			expression.Name("InherentRisk"), expression.Name("ResidualRisk"),
			expression.Name("Cost"), expression.Name("ControlRisk"),
			expression.Name("MaturityLevel"), expression.Name("ControlBusinessImpact"),
			expression.Name("ControlSecurityPriority"), expression.Name("Comment"),
			expression.Name("Attachments"),
		)
		if includeHistory {
			proj = proj.AddNames(expression.Name("HistoryData"))
		}

		keyCondition := expression.Key(ASSESSMENTID).Equal(expression.Value(assessmentId.String()))
		condition, err := expression.NewBuilder().WithKeyCondition(keyCondition).WithProjection(proj).Build()

		if err != nil {
			return entities, err
		}

		response, err := h.Repository.FindAllUsingQuery(condition, entity.TableName(), "AssessmentIdIndex", ASSESSMENTID, assessmentId.String())
		if err != nil {
			return entities, err
		}

		if response != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+"response:count::", response.Count)
			for _, value := range response.Items {
				entity, err := EntityControl.ParseDynamoAtributeToStruct(value)
				if err != nil {
					return entities, err
				}
				entities = append(entities, entity)
			}
		}
	}
	return entities, nil
}
func (h *AssesmentControlService) CreateControl(assessMentControl EntityControl.AssessmentControl, controlDataItem EntityControlData.ControlData, createdBy string) (entity EntityControl.AssessmentControl, err error) {
	const MethodName = "CreateControl"
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" assessMentControl:::", assessMentControl)
	setDefaultControlDataValues(&assessMentControl, controlDataItem, createdBy)
	_, err = h.Repository.CreateOrUpdate(assessMentControl.GetMap(), assessMentControl.TableName())
	if err != nil {
		logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:CreateOrUpdate::", err)
		return assessMentControl, err

	}
	return assessMentControl, nil
}
func (h *AssesmentControlService) DeleteAssessmentControls(assessmentId uuid.UUID) error {
	const MethodName = "DeleteAssessmentControls"
	entities, _ := h.FindAssessmentControls(assessmentId, true)

	for _, control := range entities {
		_, err := h.Repository.Delete(control.GetFilterId(), control.TableName())
		if err != nil {
			logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" err:Delete::", err)
			return err
		}

	}
	return nil
}

func setDefaultControlDataValues(assessmentControl *EntityControl.AssessmentControl, controlData EntityControlData.ControlData, createdBy string) {
	assessmentControl.ID = uuid.New()
	assessmentControl.ControlNumber = controlData.ControlNumber
	assessmentControl.ControlFamily = controlData.ControlFamily
	assessmentControl.ShortDescription = controlData.ShortDescription
	assessmentControl.MappedControl = controlData.MappedControl
	assessmentControl.Requirement = controlData.Requirement
	assessmentControl.UpdatedAt = time.Now()
	assessmentControl.CreatedAt = time.Now()
	assessmentControl.CreatedBy = createdBy
	assessmentControl.UpdatedBy = createdBy

	assessmentControl.IsEnable = false
	// status
	var status EntityControl.Map
	status.Name = "NA"
	status.Value = -1
	assessmentControl.Status = status

	assessmentControl.CurrentScore = 0.0
	assessmentControl.TargetScore = 0.0
	assessmentControl.DueDate = time.Now()
	assessmentControl.AssignedTo = "NA"

	var impact EntityControl.Map
	impact.Name = "NA"
	impact.Value = -1
	assessmentControl.Impact = impact

	var likelihood EntityControl.Map
	likelihood.Name = "NA"
	likelihood.Value = -1
	assessmentControl.Likelihood = likelihood

	assessmentControl.InherentRisk = 0
	assessmentControl.ResidualRisk = 0.0
	assessmentControl.Cost = 0
	assessmentControl.ControlRisk = "NA"
	assessmentControl.ControlBusinessImpact = "NA"
	assessmentControl.ControlSecurityPriority = "NA"
	assessmentControl.MaturityLevel = "NA"

	// comment
	var comment EntityControl.Comment
	comment.Description = "Created new assessment"
	comment.ID = uuid.New()
	comment.UpdatedAt = time.Now()
	comment.CreatedAt = time.Now()
	comment.CreatedBy = createdBy
	comment.UpdatedBy = createdBy
	comment.ControlNumber = assessmentControl.ControlNumber
	assessmentControl.Comment = comment

	//question
	questions := []EntityControl.Question{}

	var question EntityControl.Question
	question.QuestionNumber = assessmentControl.ControlNumber
	question.Description = assessmentControl.ShortDescription
	question.IsEnable = true
	question.CurrentScore = 0.0
	question.TargetScore = 0.0
	questions = append(questions, question)
	assessmentControl.Questions = questions

	//attachments
	attachments := []EntityControl.Attachment{}
	var attachment EntityControl.Attachment
	attachment.Description = "New Assessment created."
	attachment.ID = uuid.New()
	attachment.UpdatedAt = time.Now().Format(Entities.GetTimeFormat())
	attachment.CreatedAt = time.Now().Format(Entities.GetTimeFormat())
	attachment.CreatedBy = createdBy
	attachment.UpdatedBy = createdBy
	attachment.IsFile = true
	attachment.Path = " "
	attachments = append(attachments, attachment)
	assessmentControl.Attachments = attachments
	//question
	historySlice := []EntityControl.History{}
	var history EntityControl.History
	history.Description = "New Assessment created."
	history.ID = uuid.New()
	history.UpdatedAt = time.Now().Format(Entities.GetTimeFormat())
	history.CreatedAt = time.Now().Format(Entities.GetTimeFormat())
	history.CreatedBy = createdBy
	history.UpdatedBy = createdBy

	historySlice = append(historySlice, history)
	assessmentControl.HistoryData = historySlice

}
