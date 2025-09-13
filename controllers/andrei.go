package controllers

import (
	"net/http"
	"strconv"

	"andrei-api/config"
	"andrei-api/models"
	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := config.DB.Preload("Posts").Preload("Reports").Preload("Rewards").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func CreateReward(c *gin.Context) {
	var input models.RewardCreate

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var demon models.User
	if err := config.DB.Where("id = ? AND role = ?", input.DemonID, models.RoleDemon).First(&demon).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Demon not found"})
		return
	}

	reward := models.Reward{
		DemonID:     input.DemonID,
		Type:        input.Type,
		Title:       input.Title,
		Description: input.Description,
		Points:      input.Points,
	}

	if err := config.DB.Create(&reward).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reward"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"reward": reward})
}

func GetPlatformStats(c *gin.Context) {
	var stats models.PlatformStats

	config.DB.Model(&models.User{}).Count(&stats.TotalUsers)
	config.DB.Model(&models.User{}).Where("role = ?", models.RoleDemon).Count(&stats.TotalDemons)
	config.DB.Model(&models.User{}).Where("role = ?", models.RoleNetworkAdmin).Count(&stats.TotalNetworkAdmins)
	config.DB.Model(&models.Post{}).Count(&stats.TotalPosts)
	config.DB.Model(&models.Report{}).Count(&stats.TotalReports)

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func GetDemonRanking(c *gin.Context) {
	var demons []models.User
	if err := config.DB.Where("role = ?", models.RoleDemon).Find(&demons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch demons"})
		return
	}

	var demonStats []models.DemonStats
	for _, demon := range demons {
		var stats models.DemonStats
		stats.DemonID = demon.ID

		config.DB.Model(&models.User{}).Where("role = ? AND id IN (SELECT victim_id FROM reports WHERE demon_id = ?)", 
			models.RoleNetworkAdmin, demon.ID).Count(&stats.VictimsCount)
		
		config.DB.Model(&models.Reward{}).Where("demon_id = ? AND type = ?", demon.ID, models.RewardTypeReward).Count(&stats.RewardsCount)
		config.DB.Model(&models.Reward{}).Where("demon_id = ? AND type = ?", demon.ID, models.RewardTypePunishment).Count(&stats.PunishmentsCount)
		config.DB.Model(&models.Report{}).Where("demon_id = ?", demon.ID).Count(&stats.ReportsCount)
		
		config.DB.Model(&models.Reward{}).Where("demon_id = ?", demon.ID).Select("COALESCE(SUM(points), 0)").Scan(&stats.TotalPoints)

		demonStats = append(demonStats, stats)
	}

	c.JSON(http.StatusOK, gin.H{"demon_rankings": demonStats})
}

func GetAllPosts(c *gin.Context) {
	var posts []models.Post
	if err := config.DB.Preload("Author").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func DeletePost(c *gin.Context) {
	postID := c.Param("id")
	id, err := strconv.ParseUint(postID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	if err := config.DB.Delete(&models.Post{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func CreateAndreiPost(c *gin.Context) {
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