package main

import (
	"log"
	"os"

	"andrei-api/config"
	"andrei-api/models"
	"andrei-api/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func runSeeder() {
	log.Println("Running seeder...")

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

	log.Println("✅ Seeder completed successfully!")
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	config.ConnectDatabase()

	// Check if seeder flag is passed
	if len(os.Args) > 1 && os.Args[1] == "-seed" {
		runSeeder()
		return
	}

	// Create Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup routes
	routes.SetupRoutes(r)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}