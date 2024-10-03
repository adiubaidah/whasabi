package model

import (
	"time"
)

type CreateProcessModel struct {
	Name        string  `gorm:"not null" json:"name" validate:"required"`
	Phone       string  `gorm:"not null;unique" json:"phone" validate:"required,number"`
	Instruction string  `gorm:"not null" json:"instruction" validate:"required"`
	TopK        int32   `gorm:"column:top_k;type:int" json:"top_k"`
	TopP        float32 `gorm:"column:top_k;type:int" json:"top_p"`
	Temperature float32 `gorm:"type:float;not null" json:"temperature" validate:"required"`
}

type Process struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"column:user_id;not null" json:"user_id"` // Foreign key to User
	CreateProcessModel
	IsActive        bool      `gorm:"default:false" json:"is_active"`
	IsAuthenticated bool      `gorm:"default:false" json:"is_authenticated"`
	CreatedAt       time.Time `gorm:"default:autoCreateTime" json:"created_at"`
	User            *User     `gorm:"foreignKey:user_id;references:id"`
	Histories       []History `gorm:"foreignKey:service_id;references:id"`
}

func (ai *Process) TableName() string {
	return "process"
}
