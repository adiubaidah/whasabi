package ai

import "time"

type Ai struct {
	ID              uint      `gorm:"primaryKey"`
	UserID          uint      `gorm:"not null"`
	Name            string    `gorm:"not null"`
	Phone           string    `gorm:"not null;unique"`
	Instruction     string    `gorm:"not null"`
	Temperature     float64   `gorm:"type:float;not null"`
	TopK            int32     `gorm:"column:top_k;type:int"`
	TopP            float32   `gorm:"column:top_p;type:float"`
	IsActive        bool      `gorm:"default:false"`
	IsAuthenticated bool      `gorm:"default:false"`
	CreatedAt       time.Time `gorm:"default:autoCreateTime"`
}

func (ai *Ai) TableName() string {
	return "services"
}
