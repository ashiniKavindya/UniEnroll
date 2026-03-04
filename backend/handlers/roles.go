package handlers

import (
	"encoding/csv"
	"fmt"
	"io"
	"context"
	"net/http"
	"strconv"
	"strings"
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

type importStudentFailure struct {
	Row   int    `json:"row"`
	Email string `json:"email,omitempty"`
	Error string `json:"error"`
}

type importedStudentCredential struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserID   string `json:"userID"`
	Password string `json:"tempPassword"`
}

// ImportStudentsCSV imports student accounts from CSV (Admin only)
// Required headers: name,email,departmentID
// Optional headers: yearOfStudy,studentNumber,enrollmentYear
func ImportStudentsCSV(c *gin.Context) {
	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing CSV file field 'file'"})
		return
	}

	opened, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer opened.Close()

	reader := csv.NewReader(opened)
	headers, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CSV file is empty"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid CSV header"})
		return
	}

	normalize := func(s string) string {
		s = strings.TrimSpace(strings.ToLower(s))
		s = strings.ReplaceAll(s, "_", "")
		s = strings.ReplaceAll(s, "-", "")
		s = strings.ReplaceAll(s, " ", "")
		return s
	}

	idx := map[string]int{}
	for i, h := range headers {
		idx[normalize(h)] = i
	}

	find := func(keys ...string) int {
		for _, k := range keys {
			if v, exists := idx[k]; exists {
				return v
			}
		}
		return -1
	}

	nameIdx := find("name", "fullname")
	emailIdx := find("email")
	departmentIdx := find("departmentid", "department")
	yearIdx := find("yearofstudy", "year")
	studentNumberIdx := find("studentnumber")
	enrollmentYearIdx := find("enrollmentyear")

	if nameIdx == -1 || emailIdx == -1 || departmentIdx == -1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "CSV must include headers: name,email,departmentID",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	created := 0
	skipped := 0
	processed := 0
	failures := make([]importStudentFailure, 0)
	credentials := make([]importedStudentCredential, 0)

	currentYear := time.Now().Year()
	rowNumber := 1

	for {
		row, readErr := reader.Read()
		if readErr == io.EOF {
			break
		}
		rowNumber++

		if readErr != nil {
			failures = append(failures, importStudentFailure{Row: rowNumber, Error: "invalid CSV row format"})
			continue
		}

		if len(row) == 0 {
			continue
		}

		get := func(i int) string {
			if i < 0 || i >= len(row) {
				return ""
			}
			return strings.TrimSpace(row[i])
		}

		name := get(nameIdx)
		email := strings.ToLower(get(emailIdx))
		departmentID := get(departmentIdx)

		if name == "" || email == "" || departmentID == "" {
			failures = append(failures, importStudentFailure{Row: rowNumber, Email: email, Error: "name, email and departmentID are required"})
			continue
		}

		yearOfStudy := 1
		if y := get(yearIdx); y != "" {
			parsedYear, parseErr := strconv.Atoi(y)
			if parseErr != nil || parsedYear < 1 {
				failures = append(failures, importStudentFailure{Row: rowNumber, Email: email, Error: "invalid yearOfStudy"})
				continue
			}
			yearOfStudy = parsedYear
		}

		enrollmentYear := currentYear
		if ey := get(enrollmentYearIdx); ey != "" {
			parsedEnrollmentYear, parseErr := strconv.Atoi(ey)
			if parseErr != nil || parsedEnrollmentYear < 1900 {
				failures = append(failures, importStudentFailure{Row: rowNumber, Email: email, Error: "invalid enrollmentYear"})
				continue
			}
			enrollmentYear = parsedEnrollmentYear
		}

		processed++

		q := f.Collection("users").Where("email", "==", email).Limit(1)
		docs := q.Documents(ctx)
		if _, existingErr := docs.Next(); existingErr == nil {
			skipped++
			failures = append(failures, importStudentFailure{Row: rowNumber, Email: email, Error: "email already exists"})
			continue
		}

		password, genErr := utils.GeneratePassword(12)
		if genErr != nil {
			failures = append(failures, importStudentFailure{Row: rowNumber, Email: email, Error: "failed to generate password"})
			continue
		}

		hash, hashErr := bcrypt.GenerateFromPassword([]byte(password), 10)
		if hashErr != nil {
			failures = append(failures, importStudentFailure{Row: rowNumber, Email: email, Error: "failed to hash password"})
			continue
		}

		user := models.User{
			Name:                name,
			Email:               email,
			PasswordHash:        string(hash),
			Role:                models.RoleStudent,
			ForcePasswordChange: true,
		}

		userRef, _, userErr := f.Collection("users").Add(ctx, user)
		if userErr != nil {
			failures = append(failures, importStudentFailure{Row: rowNumber, Email: email, Error: "failed to create user"})
			continue
		}

		student := models.Student{
			UserID:         userRef.ID,
			StudentNumber:  get(studentNumberIdx),
			DepartmentID:   departmentID,
			YearOfStudy:    yearOfStudy,
			EnrollmentYear: enrollmentYear,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		_, studentErr := f.Collection("students").Doc(userRef.ID).Set(ctx, student)
		if studentErr != nil {
			_, _ = f.Collection("users").Doc(userRef.ID).Delete(ctx)
			failures = append(failures, importStudentFailure{Row: rowNumber, Email: email, Error: "failed to create student profile"})
			continue
		}

		credentials = append(credentials, importedStudentCredential{
			Name:     name,
			Email:    email,
			UserID:   userRef.ID,
			Password: password,
		})

		go utils.SendWelcomeEmail(email, name, userRef.ID, password, models.RoleStudent)
		created++
	}

	if created == 0 && len(failures) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":   "no students imported",
			"processed": processed,
			"created":   created,
			"skipped":   skipped,
			"failures":  failures,
		})
		return
	}

	status := http.StatusOK
	if len(failures) > 0 {
		status = http.StatusMultiStatus
	}

	c.JSON(status, gin.H{
		"message":      fmt.Sprintf("import completed: %d created, %d skipped", created, skipped),
		"processed":    processed,
		"created":      created,
		"skipped":      skipped,
		"credentials":  credentials,
		"failures":     failures,
		"requiredCSV":  []string{"name", "email", "departmentID"},
		"optionalCSV":  []string{"yearOfStudy", "studentNumber", "enrollmentYear"},
	})
}
