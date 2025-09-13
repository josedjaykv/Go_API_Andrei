package models

import (
	"time"
	"gorm.io/gorm"
)

type DemonVictim struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	DemonID   uint           `json:"demon_id" gorm:"not null"`
	Demon     User           `json:"demon" gorm:"foreignKey:DemonID"`
	VictimID  uint           `json:"victim_id" gorm:"not null"`
	Victim    User           `json:"victim" gorm:"foreignKey:VictimID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type AssignVictimRequest struct {
	VictimID uint `json:"victim_id" binding:"required"`
}