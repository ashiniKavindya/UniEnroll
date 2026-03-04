package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"backend/models"
	"backend/utils"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// retrieve firestore client from context middleware
	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("[LOGIN DEBUG] Searching for email: %s", req.Email)
	q := f.Collection("users").Where("email", "==", req.Email).Limit(1)
	docs := q.Documents(ctx)
	doc, err := docs.Next()
	if err != nil {
		log.Printf("[LOGIN DEBUG] User not found or query error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	log.Printf("[LOGIN DEBUG] User found with ID: %s", doc.Ref.ID)

	var u models.User
	if err := doc.DataTo(&u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	log.Printf("[LOGIN DEBUG] User data retrieved - Email: %s, Role: %s, PasswordHashLen: %d", u.Email, u.Role, len(u.PasswordHash))

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		log.Printf("[LOGIN DEBUG] Password comparison failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	log.Printf("[LOGIN DEBUG] Password comparison successful")

	// generate token
	token, err := utils.GenerateToken(doc.Ref.ID, u.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":               token,
		"role":                u.Role,
		"forcePasswordChange": u.ForcePasswordChange,
	})
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=admin student lecturer"`
}

func Register(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{"error": "self-registration is disabled; contact an admin"})
	return

	var req registerRequest
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

	// check if email already exists
	q := f.Collection("users").Where("email", "==", req.Email).Limit(1)
	docs := q.Documents(ctx)
	if _, err := docs.Next(); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// create user doc
	u := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         req.Role,
	}

	_, err = f.Collection("users").NewDoc().Set(ctx, u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

type changePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6"`
}

// ChangePassword allows authenticated users to update their password
func ChangePassword(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req changePasswordRequest
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

	doc, err := f.Collection("users").Doc(userID.(string)).Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var u models.User
	if err := doc.DataTo(&u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "current password is incorrect"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	_, err = f.Collection("users").Doc(userID.(string)).Update(ctx, []firestore.Update{
		{Path: "passwordHash", Value: string(hash)},
		{Path: "forcePasswordChange", Value: false},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated"})
}

// GetProfile returns the current authenticated user's profile
func GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc, err := f.Collection("users").Doc(userID.(string)).Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var u models.User
	if err := doc.DataTo(&u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	response := gin.H{
		"id":    userID,
		"name":  u.Name,
		"email": u.Email,
		"role":  role,
	}

	// If student, fetch student profile for departmentID
	if role == models.RoleStudent {
		q := f.Collection("students").Where("userID", "==", userID.(string)).Limit(1)
		docs := q.Documents(ctx)
		studentDoc, err := docs.Next()
		if err == nil {
			var student models.Student
			if err := studentDoc.DataTo(&student); err == nil {
				response["departmentID"] = student.DepartmentID
				response["yearOfStudy"] = student.YearOfStudy
			}
		}
	}

	// If lecturer, include lecturerID (which is userID)
	if role == models.RoleLecturer {
		response["lecturerID"] = userID
	}

	c.JSON(http.StatusOK, response)
}
