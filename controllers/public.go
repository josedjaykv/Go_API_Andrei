package controllers

import (
	"net/http"

	"andrei-api/config"
	"andrei-api/models"
	"github.com/gin-gonic/gin"
)

func GetResistancePage(c *gin.Context) {
	var posts []models.Post
	if err := config.DB.Preload("Author").Order("created_at DESC").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	var publicPosts []gin.H
	for _, post := range posts {
		postData := gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"body":       post.Body,
			"media":      post.Media,
			"anonymous":  post.Anonymous,
			"created_at": post.CreatedAt,
		}

		if post.Anonymous {
			postData["author"] = "Anonymous"
		} else if post.Author != nil {
			postData["author"] = post.Author.Username
		} else {
			postData["author"] = "Unknown"
		}

		publicPosts = append(publicPosts, postData)
	}

	c.JSON(http.StatusOK, gin.H{"posts": publicPosts})
}