package models

import (
	"time"
	"gorm.io/gorm"
)

type Report struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	DemonID     uint           `json:"demon_id" gorm:"not null"`
	Demon       User           `json:"demon" gorm:"foreignKey:DemonID"`
	VictimID    uint           `json:"victim_id" gorm:"not null"`
	Victim      User           `json:"victim" gorm:"foreignKey:VictimID"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"not null"`
	Status      string         `json:"status" gorm:"default:'pending'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type ReportCreate struct {
	VictimID    uint   `json:"victim_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}