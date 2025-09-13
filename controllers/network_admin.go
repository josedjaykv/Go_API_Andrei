package controllers

import (
	"net/http"

	"andrei-api/config"
	"andrei-api/models"
	"github.com/gin-gonic/gin"
)

func CreateAnonymousPost(c *gin.Context) {
	var input models.PostCreate

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := models.Post{
		Title:     input.Title,
		Body:      input.Body,
		Media:     input.Media,
		AuthorID:  nil,
		Anonymous: true,
	}

	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"post": post})
}