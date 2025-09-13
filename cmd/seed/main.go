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

	// Check if Andrei user exists
	var existingUser models.User
	err := config.DB.Where("email = ?", "andrei@evil.com").First(&existingUser).Error
	if err == nil {
		log.Println("✅ Andrei user already exists")
		return
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatal("❌ Error checking for Andrei user:", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("❌ Failed to hash password:", err)
	}

	// Create Andrei user
	andrei := models.User{
		Username: "AndreiMesManur",
		Email:    "andrei@evil.com",
		Password: string(hashedPassword),
		Role:     models.RoleAndrei,
	}

	if err := config.DB.Create(&andrei).Error; err != nil {
		log.Fatal("❌ Failed to create Andrei user:", err)
	}

	log.Println("🎉 Andrei user created successfully!")
	log.Println("   Email:    andrei@evil.com")
	log.Println("   Password: password123")
}
