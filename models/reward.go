package models

import (
	"time"
	"gorm.io/gorm"
)

type RewardType string

const (
	RewardTypeReward     RewardType = "reward"
	RewardTypePunishment RewardType = "punishment"
)

type Reward struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	DemonID     uint           `json:"demon_id" gorm:"not null"`
	Demon       User           `json:"demon" gorm:"foreignKey:DemonID"`
	Type        RewardType     `json:"type" gorm:"not null"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"not null"`
	Points      int            `json:"points" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type RewardCreate struct {
	DemonID     uint       `json:"demon_id" binding:"required"`
	Type        RewardType `json:"type" binding:"required"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description" binding:"required"`
	Points      int        `json:"points"`
}