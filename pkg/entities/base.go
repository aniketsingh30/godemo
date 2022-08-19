package entities

import (
	"time"

	"github.com/google/uuid"
)

type Interface interface {
	GenerateID()
	SetCreatedAt()
	SetUpdatedAt()
	TableName() string
	GetMap() map[string]interface{}
	GetFilterId() map[string]interface{}
}

type Base struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedBy string    `json:"createdBy"`
	UpdatedBy string    `json:"updatedBy"`
}

func (b *Base) GenerateID() {
	b.ID = uuid.New()
}

func (b *Base) SetCreatedAt() {
	b.CreatedAt = time.Now()
}

func (b *Base) SetUpdatedAt() {
	b.UpdatedAt = time.Now()
}

func GetTimeFormat() string {
	return "2006-01-02T15:04:05"
}
