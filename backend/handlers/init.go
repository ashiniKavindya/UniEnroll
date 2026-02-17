package handlers

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"backend/models"
)

// InitDB initializes the database with sample data
func InitDB(c *gin.Context) {
	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create sample users
	adminHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10)
	studentHash, _ := bcrypt.GenerateFromPassword([]byte("student123"), 10)

	users := map[string]models.User{
		"admin": {
			Name:         "Admin User",
			Email:        "admin@uni.edu",
			PasswordHash: string(adminHash),
			Role:         "admin",
		},
		"student1": {
			Name:         "John Doe",
			Email:        "student1@uni.edu",
			PasswordHash: string(studentHash),
			Role:         "user",
		},
		"student2": {
			Name:         "Jane Smith",
			Email:        "student2@uni.edu",
			PasswordHash: string(studentHash),
			Role:         "user",
		},
	}

	batch := f.Batch()
	for _, u := range users {
		doc := f.Collection("users").NewDoc()
		batch.Set(doc, u)
	}

	// Create sample courses
	courses := map[string]models.Course{
		"cs101": {
			Title:       "Introduction to Computer Science",
			Code:        "CS101",
			Description: "Fundamentals of programming and algorithms",
			Instructor:  "Dr. Smith",
			Credits:     3,
		},
		"math201": {
			Title:       "Calculus II",
			Code:        "MATH201",
			Description: "Advanced calculus concepts",
			Instructor:  "Dr. Johnson",
			Credits:     4,
		},
		"eng102": {
			Title:       "English Literature",
			Code:        "ENG102",
			Description: "Classic and contemporary literature",
			Instructor:  "Prof. Williams",
			Credits:     3,
		},
	}

	for _, c := range courses {
		doc := f.Collection("courses").NewDoc()
		batch.Set(doc, c)
	}

	// Create sample modules
	modules := []models.Module{
		{
			Title:       "Variables and Data Types",
			CourseID:    "cs101",
			Description: "Learn about basic programming concepts",
			Credits:     1,
		},
		{
			Title:       "Functions and Loops",
			CourseID:    "cs101",
			Description: "Control flow and reusable code",
			Credits:     1,
		},
		{
			Title:       "Derivatives",
			CourseID:    "math201",
			Description: "Understanding rates of change",
			Credits:     2,
		},
	}

	for _, m := range modules {
		doc := f.Collection("modules").NewDoc()
		batch.Set(doc, m)
	}

	_, err := batch.Commit(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initialize database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Database initialized successfully",
		"data": gin.H{
			"users":   len(users),
			"courses": len(courses),
			"modules": len(modules),
		},
	})
}
