package models

import (
	"time"
	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"not null"`
	Body      string         `json:"body" gorm:"not null"`
	Media     string         `json:"media,omitempty"`
	AuthorID  *uint          `json:"author_id,omitempty"`
	Author    *User          `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Anonymous bool           `json:"anonymous" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type PostCreate struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
	Media string `json:"media,omitempty"`
}