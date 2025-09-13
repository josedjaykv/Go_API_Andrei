package routes

import (
	"andrei-api/controllers"
	"andrei-api/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	// Public routes
	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)
	api.GET("/resistance", controllers.GetResistancePage)

	// Protected routes
	auth := api.Group("/")
	auth.Use(middleware.AuthRequired())

	// Andrei routes (admin only)
	andrei := auth.Group("/admin")
	andrei.Use(middleware.RequireAndrei())
	{
		andrei.GET("/users", controllers.GetAllUsers)
		andrei.GET("/users/:id", controllers.GetUserByID)
		andrei.DELETE("/users/:id", controllers.DeleteUser)
		andrei.POST("/rewards", controllers.CreateReward)
		andrei.GET("/stats", controllers.GetPlatformStats)
		andrei.GET("/demons/ranking", controllers.GetDemonRanking)
		andrei.GET("/posts", controllers.GetAllPosts)
		andrei.DELETE("/posts/:id", controllers.DeletePost)
		andrei.POST("/posts", controllers.CreateAndreiPost)
	}

	// Demon routes
	demons := auth.Group("/demons")
	demons.Use(middleware.RequireDemon())
	{
		demons.POST("/victims", controllers.RegisterVictim)
		demons.POST("/reports", controllers.CreateReport)
		demons.GET("/stats", controllers.GetMyStats)
		demons.GET("/victims", controllers.GetMyVictims)
		demons.GET("/reports", controllers.GetMyReports)
		demons.PUT("/reports/:id", controllers.UpdateReportStatus)
		demons.POST("/posts", controllers.CreateDemonPost)
	}

	// Network Admin routes
	networkAdmins := auth.Group("/network-admins")
	networkAdmins.Use(middleware.RequireNetworkAdmin())
	{
		networkAdmins.POST("/posts/anonymous", controllers.CreateAnonymousPost)
	}

	// Mixed role routes (Andrei and Demons can create posts)
	posts := auth.Group("/posts")
	posts.Use(middleware.RequireAndreiOrDemon())
	{
		// These are already handled in specific role routes above
	}
}