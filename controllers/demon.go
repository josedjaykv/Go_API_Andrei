package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"andrei-api/config"
	"andrei-api/models"

	"github.com/gin-gonic/gin"
)

func GetAvailableNetworkAdmins(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	// Obtener todos los network admins que no son víctimas de este demonio
	var networkAdmins []models.User
	if err := config.DB.Where("role = ? AND id NOT IN (SELECT victim_id FROM demon_victims WHERE demon_id = ?)",
		models.RoleNetworkAdmin, user.ID).Find(&networkAdmins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available network admins"})
		return
	}

	// Limpiar los datos para evitar problemas con campos NULL
	var cleanNetworkAdmins []gin.H
	for _, admin := range networkAdmins {
		cleanNetworkAdmins = append(cleanNetworkAdmins, gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"email":    admin.Email,
			"role":     admin.Role,
		})
	}

	c.JSON(http.StatusOK, gin.H{"available_network_admins": cleanNetworkAdmins})
}

func AssignVictim(c *gin.Context) {
	var input models.AssignVictimRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(models.User)

	// Verificar que la víctima existe y es network admin
	var victim models.User
	if err := config.DB.Where("id = ? AND role = ?", input.VictimID, models.RoleNetworkAdmin).First(&victim).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Network admin not found"})
		return
	}

	// Verificar que no sea ya víctima de este demonio
	var existingRelation models.DemonVictim
	if err := config.DB.Where("demon_id = ? AND victim_id = ?", user.ID, input.VictimID).First(&existingRelation).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "This network admin is already your victim"})
		return
	}

	// Crear la relación demonio-víctima
	demonVictim := models.DemonVictim{
		DemonID:  user.ID,
		VictimID: input.VictimID,
	}

	if err := config.DB.Create(&demonVictim).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign victim"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Victim assigned successfully",
		"victim": gin.H{
			"id":       victim.ID,
			"username": victim.Username,
			"email":    victim.Email,
			"role":     victim.Role,
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

	var demonVictims []models.DemonVictim
	if err := config.DB.Where("demon_id = ?", user.ID).Preload("Victim").Find(&demonVictims).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch victims"})
		return
	}

	// Extraer solo la información de las víctimas
	var victims []models.User
	for _, dv := range demonVictims {
		victims = append(victims, dv.Victim)
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
