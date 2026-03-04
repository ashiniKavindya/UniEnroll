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

	batch := f.Batch()

	// Create sample faculties
	facultyScience := f.Collection("faculties").NewDoc()
	facultyEng := f.Collection("faculties").NewDoc()
	facultyArts := f.Collection("faculties").NewDoc()

	faculties := map[*firestore.DocumentRef]models.Faculty{
		facultyScience: {
			FacultyID:   facultyScience.ID,
			Name:        "Faculty of Science",
			Code:        "SCI",
			Description: "Natural sciences and mathematics",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		facultyEng: {
			FacultyID:   facultyEng.ID,
			Name:        "Faculty of Engineering",
			Code:        "ENG",
			Description: "Engineering and technology",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		facultyArts: {
			FacultyID:   facultyArts.ID,
			Name:        "Faculty of Arts",
			Code:        "ARTS",
			Description: "Humanities and social sciences",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for doc, faculty := range faculties {
		batch.Set(doc, faculty)
	}

	// Create sample departments
	deptCS := f.Collection("departments").NewDoc()
	deptMath := f.Collection("departments").NewDoc()
	deptEng := f.Collection("departments").NewDoc()

	departments := map[*firestore.DocumentRef]models.Department{
		deptCS: {
			DepartmentID: deptCS.ID,
			Name:         "Department of Computer Science",
			Code:         "CS",
			FacultyID:    facultyScience.ID,
			Description:  "Computing and software development",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		deptMath: {
			DepartmentID: deptMath.ID,
			Name:         "Department of Mathematics",
			Code:         "MATH",
			FacultyID:    facultyScience.ID,
			Description:  "Pure and applied mathematics",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		deptEng: {
			DepartmentID: deptEng.ID,
			Name:         "Department of English",
			Code:         "ENG",
			FacultyID:    facultyArts.ID,
			Description:  "English language and literature",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	for doc, dept := range departments {
		batch.Set(doc, dept)
	}

	adminHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10)
	lecturerHash, _ := bcrypt.GenerateFromPassword([]byte("lecturer123"), 10)
	studentHash, _ := bcrypt.GenerateFromPassword([]byte("student123"), 10)

	users := map[string]models.User{
		"admin": {
			Name:         "Admin User",
			Email:        "admin@uni.edu",
			PasswordHash: string(adminHash),
			Role:         "admin",
			ForcePasswordChange: false,
		},
		"lecturer": {
			Name:         "John Lecturer",
			Email:        "john@uni.edu",
			PasswordHash: string(lecturerHash),
			Role:         "lecturer",
			ForcePasswordChange: false,
		},
		"student": {
			Name:         "Jane Student",
			Email:        "jane@uni.edu",
			PasswordHash: string(studentHash),
			Role:         "student",
			ForcePasswordChange: false,
		},
	}

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
			"faculties":   len(faculties),
			"departments": len(departments),
			"users":       len(users),
			"courses":     len(courses),
			"modules":     len(modules),
		},
	})
}
