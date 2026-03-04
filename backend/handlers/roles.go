package handlers

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"backend/models"
	"backend/utils"
)

// CreateStudentProfile creates a student profile when a student user registers
func CreateStudentProfile(c *gin.Context) {
	type studentRequest struct {
		UserID         string `json:"userID" binding:"required"`
		StudentNumber  string `json:"studentNumber"`
		DepartmentID   string `json:"departmentID" binding:"required"`
		YearOfStudy    int    `json:"yearOfStudy"`
		EnrollmentYear int    `json:"enrollmentYear"`
	}

	var req studentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	student := models.Student{
		UserID:         req.UserID,
		StudentNumber:  req.StudentNumber,
		DepartmentID:   req.DepartmentID,
		YearOfStudy:    req.YearOfStudy,
		EnrollmentYear: req.EnrollmentYear,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, err := f.Collection("students").Doc(req.UserID).Set(ctx, student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create student profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "student profile created", "userID": req.UserID})
}

// CreateLecturerProfile creates a lecturer profile
func CreateLecturerProfile(c *gin.Context) {
	type lecturerRequest struct {
		UserID         string   `json:"userID" binding:"required"`
		Department     string   `json:"department"`
		Specialization string   `json:"specialization"`
		OfficeLocation string   `json:"officeLocation"`
		PhoneNumber    string   `json:"phoneNumber"`
		Qualifications []string `json:"qualifications"`
	}

	var req lecturerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	lecturer := models.Lecturer{
		UserID:         req.UserID,
		Department:     req.Department,
		Specialization: req.Specialization,
		OfficeLocation: req.OfficeLocation,
		PhoneNumber:    req.PhoneNumber,
		Qualifications: req.Qualifications,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, err := f.Collection("lecturers").Doc(req.UserID).Set(ctx, lecturer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create lecturer profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "lecturer profile created", "userID": req.UserID})
}

// CreateAdminProfile creates an admin profile
func CreateAdminProfile(c *gin.Context) {
	type adminRequest struct {
		UserID      string   `json:"userID" binding:"required"`
		Permissions []string `json:"permissions"`
		Department  string   `json:"department"`
		PhoneNumber string   `json:"phoneNumber"`
	}

	var req adminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	admin := models.Admin{
		UserID:      req.UserID,
		Permissions: req.Permissions,
		Department:  req.Department,
		PhoneNumber: req.PhoneNumber,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := f.Collection("admins").Doc(req.UserID).Set(ctx, admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create admin profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "admin profile created", "userID": req.UserID})
}

// CreateModule creates a new module/subject
func CreateModule(c *gin.Context) {
	type moduleRequest struct {
		Title       string `json:"title" binding:"required"`
		Code        string `json:"code" binding:"required"`
		CourseID    string `json:"courseID" binding:"required"`
		LecturerID  string `json:"lecturerID" binding:"required"`
		Description string `json:"description"`
		Credits     int    `json:"credits"`
		Semester    int    `json:"semester"`
	}

	var req moduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create new document with auto-generated ID
	docRef, _, err := f.Collection("modules").Add(ctx, map[string]interface{}{
		"title":       req.Title,
		"code":        req.Code,
		"courseID":    req.CourseID,
		"lecturerID":  req.LecturerID,
		"description": req.Description,
		"credits":     req.Credits,
		"semester":    req.Semester,
		"createdAt":   time.Now(),
		"updatedAt":   time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create module"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "module created", "moduleID": docRef.ID})
}

// GetModule retrieves a module by ID
func GetModule(c *gin.Context) {
	moduleID := c.Param("moduleID")

	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc, err := f.Collection("modules").Doc(moduleID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "module not found"})
		return
	}

	var module models.Module
	if err := doc.DataTo(&module); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, module)
}

// GetAllModules retrieves all modules (filtered by lecturerID for lecturers)
func GetAllModules(c *gin.Context) {
	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get user role and ID from context
	role, _ := c.Get("role")
	userID, _ := c.Get("userID")

	var docs []*firestore.DocumentSnapshot
	var err error

	// If lecturer, filter by lecturerID
	if role == models.RoleLecturer {
		query := f.Collection("modules").Where("lecturerID", "==", userID.(string))
		docs, err = query.Documents(ctx).GetAll()
	} else {
		// For students and admins, show all modules
		docs, err = f.Collection("modules").Documents(ctx).GetAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch modules"})
		return
	}

	var modules []models.Module
	for _, doc := range docs {
		var module models.Module
		if err := doc.DataTo(&module); err != nil {
			continue
		}
		module.ModuleID = doc.Ref.ID
		modules = append(modules, module)
	}

	c.JSON(http.StatusOK, gin.H{"modules": modules})
}

// CreateUserAccount creates a complete user account with auto-generated password (Admin only)
func CreateUserAccount(c *gin.Context) {
	type userCreateRequest struct {
		Email          string `json:"email" binding:"required,email"`
		Name           string `json:"name" binding:"required"`
		Role           string `json:"role" binding:"required,oneof=student lecturer"`
		DepartmentID   string `json:"departmentID"`
		YearOfStudy    int    `json:"yearOfStudy"`
		StudentNumber  string `json:"studentNumber"`
		Specialization string `json:"specialization"`
		OfficeLocation string `json:"officeLocation"`
	}

	var req userCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if email already exists
	q := f.Collection("users").Where("email", "==", req.Email).Limit(1)
	docs := q.Documents(ctx)
	if _, err := docs.Next(); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	// Generate random password
	password, err := utils.GeneratePassword(12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate password"})
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Create user document
	user := models.User{
		Name:                req.Name,
		Email:               req.Email,
		PasswordHash:        string(hash),
		Role:                req.Role,
		ForcePasswordChange: true,
	}

	userRef, _, err := f.Collection("users").Add(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	userID := userRef.ID

	// Create role-specific profile
	if req.Role == models.RoleStudent {
		student := models.Student{
			UserID:         userID,
			StudentNumber:  req.StudentNumber,
			DepartmentID:   req.DepartmentID,
			YearOfStudy:    req.YearOfStudy,
			EnrollmentYear: time.Now().Year(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		_, err = f.Collection("students").Doc(userID).Set(ctx, student)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create student profile"})
			return
		}
	} else if req.Role == models.RoleLecturer {
		lecturer := models.Lecturer{
			UserID:         userID,
			Department:     req.DepartmentID,
			Specialization: req.Specialization,
			OfficeLocation: req.OfficeLocation,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		_, err = f.Collection("lecturers").Doc(userID).Set(ctx, lecturer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create lecturer profile"})
			return
		}
	}

	// Send welcome email with credentials
	go utils.SendWelcomeEmail(req.Email, req.Name, userID, password, req.Role)

	c.JSON(http.StatusCreated, gin.H{
		"message":      "user created successfully and email sent",
		"userID":       userID,
		"email":        req.Email,
		"tempPassword": password, // Return in response for admin reference (in dev)
	})
}
