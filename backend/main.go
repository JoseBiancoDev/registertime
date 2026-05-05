package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bianquiviri/control-horario/handlers"
	"github.com/bianquiviri/control-horario/middleware"
	"github.com/bianquiviri/control-horario/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database
	utils.InitDB()

	r := gin.Default()

	// Security & Rate Limiting
	r.Use(middleware.SecurityMiddleware())

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Update this for production
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Auth & Password routes
	r.POST("/api/login", handlers.Login)
	r.POST("/api/forgot-password", handlers.ForgotPassword)
	r.POST("/api/reset-password", handlers.ResetPassword)

	// Protected routes
	authorized := r.Group("/api")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("/logs", handlers.GetLogs)
		authorized.POST("/logs/start", handlers.StartLog)
		authorized.POST("/logs/stop", handlers.StopLog)
		authorized.GET("/report", handlers.GenerateExcelReport)
		authorized.POST("/change-password", handlers.ChangePassword)

		// Admin routes
		admin := authorized.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/users", handlers.GetUsers)
			admin.POST("/users", handlers.CreateUser)
			admin.GET("/users/:id/logs", handlers.GetUserLogs)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
