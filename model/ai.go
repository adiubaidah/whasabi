package model

import (
	"time"
)

type Ai struct {
	ID              uint      `gorm:"primaryKey"`
	UserID          uint      `gorm:"column:user_id;not null"`
	Name            string    `gorm:"not null"`
	Phone           string    `gorm:"not null;unique"`
	Instruction     string    `gorm:"not null"`
	Temperature     float32   `gorm:"type:float;not null"`
	TopK            int32     `gorm:"column:top_k;type:int"`
	TopP            float32   `gorm:"column:top_p;type:float"`
	IsActive        bool      `gorm:"default:false"`
	IsAuthenticated bool      `gorm:"default:false"`
	CreatedAt       time.Time `gorm:"default:autoCreateTime"`
	User            *User     `gorm:"foreignKey:user_id;references:id"`
	Histories       []History `gorm:"foreignKey:service_id;references:id"`
}

func (ai *Ai) TableName() string {
	return "services"
}

type CreateAIModel struct {
	Name        string  `json:"name" validate:"required"`
	Phone       string  `json:"phone" validate:"required,number"`
	Instruction string  `json:"instruction" validate:"required"`
	TopK        int32   `json:"topK" validate:"required,number"`
	TopP        float32 `json:"topP" validate:"required,number"`
	Temperature float32 `json:"temperature" validate:"required"`
}
