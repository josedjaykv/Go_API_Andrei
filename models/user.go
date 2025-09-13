package models

import (
	"time"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAndrei      UserRole = "andrei"
	RoleDemon       UserRole = "demon"
	RoleNetworkAdmin UserRole = "network_admin"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique;not null"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Role      UserRole       `json:"role" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	Posts     []Post     `json:"posts,omitempty" gorm:"foreignKey:AuthorID"`
	Reports   []Report   `json:"reports,omitempty" gorm:"foreignKey:DemonID"`
	Rewards   []Reward   `json:"rewards,omitempty" gorm:"foreignKey:DemonID"`
}

type UserRegistration struct {
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	Role     UserRole `json:"role" binding:"required"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}