package main

import (
	"log"

	"andrei-api/config"
	"andrei-api/models"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables desde la raíz
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️ No .env file found (usando variables del sistema)")
	}

	// Connect to database
	config.ConnectDatabase()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("❌ Failed to hash password:", err)
	}

	// Create users if they don't exist
	users := []models.User{
		{
			Username: "AndreiMesManur",
			Email:    "andrei@evil.com",
			Password: string(hashedPassword),
			Role:     models.RoleAndrei,
		},
		{
			Username: "demon1",
			Email:    "demon1@evil.com",
			Password: string(hashedPassword),
			Role:     models.RoleDemon,
		},
		{
			Username: "admin1",
			Email:    "admin1@network.com",
			Password: string(hashedPassword),
			Role:     models.RoleNetworkAdmin,
		},
		{
			Username: "admin2",
			Email:    "admin2@network.com",
			Password: string(hashedPassword),
			Role:     models.RoleNetworkAdmin,
		},
		{
			Username: "admin3",
			Email:    "admin3@network.com",
			Password: string(hashedPassword),
			Role:     models.RoleNetworkAdmin,
		},
	}

	for _, user := range users {
		var existingUser models.User
		err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error
		if err == nil {
			log.Printf("✅ User %s already exists\n", user.Username)
			continue
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Fatal("❌ Error checking for user:", err)
		}

		if err := config.DB.Create(&user).Error; err != nil {
			log.Fatal("❌ Failed to create user:", err)
		}

		log.Printf("🎉 User %s created successfully! (Email: %s, Role: %s)\n", 
			user.Username, user.Email, user.Role)
	}
}
