package controllers

import (
	"net/http"
	"strconv"

	"andrei-api/config"
	"andrei-api/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterVictim(c *gin.Context) {
	var input models.UserRegistration

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Role != models.RoleNetworkAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only network_admin role can be registered as victim"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Victim registered successfully",
		"victim": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func CreateReport(c *gin.Context) {
	var input models.ReportCreate

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var victim models.User
	if err := config.DB.Where("id = ? AND role = ?", input.VictimID, models.RoleNetworkAdmin).First(&victim).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Victim not found"})
		return
	}

	user := c.MustGet("user").(models.User)

	report := models.Report{
		DemonID:     user.ID,
		VictimID:    input.VictimID,
		Title:       input.Title,
		Description: input.Description,
		Status:      "pending",
	}

	if err := config.DB.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create report"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"report": report})
}

func GetMyStats(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var stats models.DemonStats
	stats.DemonID = user.ID

	config.DB.Model(&models.User{}).Where("role = ? AND id IN (SELECT victim_id FROM reports WHERE demon_id = ?)", 
		models.RoleNetworkAdmin, user.ID).Count(&stats.VictimsCount)
	
	config.DB.Model(&models.Reward{}).Where("demon_id = ? AND type = ?", user.ID, models.RewardTypeReward).Count(&stats.RewardsCount)
	config.DB.Model(&models.Reward{}).Where("demon_id = ? AND type = ?", user.ID, models.RewardTypePunishment).Count(&stats.PunishmentsCount)
	config.DB.Model(&models.Report{}).Where("demon_id = ?", user.ID).Count(&stats.ReportsCount)
	
	config.DB.Model(&models.Reward{}).Where("demon_id = ?", user.ID).Select("COALESCE(SUM(points), 0)").Scan(&stats.TotalPoints)

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func GetMyVictims(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var victims []models.User
	if err := config.DB.Where("role = ? AND id IN (SELECT victim_id FROM reports WHERE demon_id = ?)", 
		models.RoleNetworkAdmin, user.ID).Find(&victims).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch victims"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"victims": victims})
}

func GetMyReports(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var reports []models.Report
	if err := config.DB.Where("demon_id = ?", user.ID).Preload("Victim").Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reports"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reports": reports})
}

func CreateDemonPost(c *gin.Context) {
	var input models.PostCreate

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(models.User)

	post := models.Post{
		Title:     input.Title,
		Body:      input.Body,
		Media:     input.Media,
		AuthorID:  &user.ID,
		Anonymous: false,
	}

	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"post": post})
}

func UpdateReportStatus(c *gin.Context) {
	reportID := c.Param("id")
	id, err := strconv.ParseUint(reportID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(models.User)

	var report models.Report
	if err := config.DB.Where("id = ? AND demon_id = ?", id, user.ID).First(&report).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	report.Status = input.Status
	if err := config.DB.Save(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"report": report})
}