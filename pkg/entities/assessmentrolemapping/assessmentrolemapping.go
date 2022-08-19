package assessmentrolemapping

import (
	"encoding/json"
	"os"

	//"github.com/google/uuid"

	"riscvue.com/pkg/utils/logger"
)

type AssessmentUserMappingRole struct {
	RoleId string `json:"roleId"`
	UserId string `json:"userId"`
}

func InterfaceToModel(data interface{}) (instance *AssessmentUserMappingRole, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}

	return instance, json.Unmarshal(bytes, &instance)
}

//func (p *AssesmentUserMappingRole) GetFilterId() map[string]interface{} {
//	return map[string]interface{}{"id": p.ID.String()}
//}

func (p *AssessmentUserMappingRole) TableName() string {
	return os.Getenv("STORAGE_ASSESMENT_ROLE_MAPPING_TABLE_NAME") + "-" + os.Getenv("ENV")
}

func (p *AssessmentUserMappingRole) Bytes() ([]byte, error) {
	return json.Marshal(p)
}

func (p *AssessmentUserMappingRole) GetMap() map[string]interface{} {
	logger.INFO("className=AssessmentUserMappingRoleEntitiy MethodName=GetMap start:::", nil)
	return map[string]interface{}{
		//	"id":          p.ID.String(),
		"RoleId": p.RoleId,
		"UserId": p.UserId,
		//	"createdAt":   p.CreatedAt.Format(entities.GetTimeFormat()),
		//	"updatedAt":   p.UpdatedAt.Format(entities.GetTimeFormat()),
		//	"createdBy":   p.CreatedBy,
		//	"updatedBy":   p.CreatedBy,
	}
}

// func ParseDynamoAtributeToStruct(response map[string]*dynamodb.AttributeValue) (p AssesmentUserMappingRole, err error) {
// 	if response == nil || (response != nil && len(response) == 0) {
// 		return p, errors.New("Item not found")
// 	}
// 	for key, value := range response {

// 		if key == "id" {
// 			p.ID, err = uuid.Parse(*value.S)
// 			if p.ID == uuid.Nil {
// 				err = errors.New("Item not found")
// 			}
// 		}
// 		if key == "RoleId" {
// 			p.RoleId = *value.S
// 		}
// 		if key == "UserId" {
// 			p.UserId = *value.S
// 		}
// 		if key == "AssesmentId" {
// 			p.AssesmentId = *value.S
// 		}
// 		if key == "createdBy" {
// 			p.CreatedBy = *value.S
// 		}
// 		if key == "updatedBy" {
// 			p.UpdatedBy = *value.S
// 		}
// 		if key == "createdAt" {
// 			p.CreatedAt, err = time.Parse(entities.GetTimeFormat(), *value.S)
// 		}
// 		if key == "updatedAt" {
// 			p.UpdatedAt, err = time.Parse(entities.GetTimeFormat(), *value.S)
// 		}
// 		if err != nil {
// 			return p, err
// 		}
// 	}

// 	return p, nil
// }
