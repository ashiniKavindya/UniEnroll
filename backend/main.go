package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"backend/config"
	"backend/handlers"
	"backend/middleware"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	ctx := context.Background()
	client, err := config.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	fmt.Println("✅ Firestore connected successfully!")

	// Setup Gin
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Make client available to handlers via middleware
	r.Use(func(c *gin.Context) {
		c.Set("firestore", client)
		c.Next()
	})

	auth := r.Group("/auth")
	{
		auth.POST("/login", handlers.Login)
		auth.POST("/register", handlers.Register)
	}

	// Protected routes (require authentication)
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", handlers.GetProfile)
		protected.POST("/change-password", handlers.ChangePassword)

		// Available to all authenticated users
		protected.GET("/modules", handlers.GetAllModules)
		protected.GET("/modules/:moduleID", handlers.GetModule)
		protected.GET("/modules/search", handlers.GetModulesByDepartment)
		protected.GET("/departments", handlers.GetAllDepartments)

		// Student routes
		protected.POST("/enrollments", handlers.EnrollModule)
		protected.GET("/enrollments/:studentID", handlers.GetStudentEnrollments)
		protected.PUT("/enrollments/:enrollmentID/drop", handlers.DropModule)

		// Admin only routes
		admin := protected.Group("")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.POST("/admin/create", handlers.CreateAdminProfile)
			admin.POST("/admin/users", handlers.CreateUserAccount)
			admin.POST("/lecturer/create", handlers.CreateLecturerProfile)
			admin.POST("/student/create", handlers.CreateStudentProfile)
			admin.POST("/modules", handlers.CreateModule)
			admin.POST("/departments", handlers.CreateDepartment)
		}
	}

	// Database initialization (run once to seed data)
	r.POST("/init", handlers.InitDB)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
