package handlers

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"

	"backend/models"
)

// CreateDepartment creates a new department
func CreateDepartment(c *gin.Context) {
	type deptRequest struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		FacultyName string `json:"facultyName"`
		HOD         string `json:"hod"`
		Description string `json:"description"`
	}

	var req deptRequest
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

	docRef, _, err := f.Collection("departments").Add(ctx, map[string]interface{}{
		"name":        req.Name,
		"code":        req.Code,
		"facultyName": req.FacultyName,
		"hod":         req.HOD,
		"description": req.Description,
		"createdAt":   time.Now(),
		"updatedAt":   time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create department"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "department created", "departmentID": docRef.ID})
}

// GetAllDepartments retrieves all departments
func GetAllDepartments(c *gin.Context) {
	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	docs, err := f.Collection("departments").Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch departments"})
		return
	}

	var departments []models.Department
	for _, doc := range docs {
		var dept models.Department
		if err := doc.DataTo(&dept); err != nil {
			continue
		}
		dept.DepartmentID = doc.Ref.ID
		departments = append(departments, dept)
	}

	c.JSON(http.StatusOK, gin.H{"departments": departments})
}

// EnrollModule enrolls a student in a module
func EnrollModule(c *gin.Context) {
	type enrollRequest struct {
		StudentID string `json:"studentID" binding:"required"`
		ModuleID  string `json:"moduleID" binding:"required"`
	}

	var req enrollRequest
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

	// Check if already enrolled
	q := f.Collection("enrollments").
		Where("studentID", "==", req.StudentID).
		Where("moduleID", "==", req.ModuleID).
		Limit(1)
	docs, err := q.Documents(ctx).GetAll()
	if err == nil && len(docs) > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "already enrolled in this module"})
		return
	}

	// Create enrollment
	enrollment := models.Enrollment{
		StudentID:  req.StudentID,
		ModuleID:   req.ModuleID,
		Status:     models.EnrollmentStatusEnrolled,
		EnrolledAt: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	docRef, _, err := f.Collection("enrollments").Add(ctx, enrollment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enroll in module"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "enrolled successfully", "enrollmentID": docRef.ID})
}

// GetStudentEnrollments gets all enrollments for a student
func GetStudentEnrollments(c *gin.Context) {
	studentID := c.Param("studentID")

	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	docs, err := f.Collection("enrollments").Where("studentID", "==", studentID).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch enrollments"})
		return
	}

	var enrollments []models.Enrollment
	for _, doc := range docs {
		var enrollment models.Enrollment
		if err := doc.DataTo(&enrollment); err != nil {
			continue
		}
		enrollment.EnrollmentID = doc.Ref.ID
		enrollments = append(enrollments, enrollment)
	}

	c.JSON(http.StatusOK, gin.H{"enrollments": enrollments})
}

// GetModulesByDepartment gets modules filtered by department and type
func GetModulesByDepartment(c *gin.Context) {
	departmentID := c.Query("departmentID")
	moduleType := c.Query("type") // compulsory, elective, faculty

	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query firestore.Query
	query = f.Collection("modules").Query
	if departmentID != "" {
		query = query.Where("departmentID", "==", departmentID)
	}
	if moduleType != "" {
		query = query.Where("type", "==", moduleType)
	}

	docs, err := query.Documents(ctx).GetAll()
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

// DropModule allows a student to drop an enrolled module
func DropModule(c *gin.Context) {
	enrollmentID := c.Param("enrollmentID")

	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := f.Collection("enrollments").Doc(enrollmentID).Update(ctx, []firestore.Update{
		{Path: "status", Value: models.EnrollmentStatusDropped},
		{Path: "updatedAt", Value: time.Now()},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to drop module"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "module dropped successfully"})
}
