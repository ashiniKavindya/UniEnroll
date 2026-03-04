package handlers

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"

	"backend/models"
)

// CreateFaculty creates a new faculty
func CreateFaculty(c *gin.Context) {
	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Dean        string `json:"dean"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create new faculty document
	doc := f.Collection("faculties").NewDoc()
	now := time.Now()

	faculty := models.Faculty{
		FacultyID:   doc.ID,
		Name:        req.Name,
		Code:        req.Code,
		Dean:        req.Dean,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err := doc.Set(ctx, faculty)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create faculty"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Faculty created successfully",
		"faculty": faculty,
	})
}

// GetAllFaculties retrieves all faculties
func GetAllFaculties(c *gin.Context) {
	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	iter := f.Collection("faculties").Documents(ctx)
	defer iter.Stop()

	var faculties []models.Faculty
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch faculties"})
			return
		}

		var faculty models.Faculty
		if err := doc.DataTo(&faculty); err != nil {
			continue
		}
		faculties = append(faculties, faculty)
	}

	c.JSON(http.StatusOK, gin.H{
		"faculties": faculties,
	})
}

// GetFaculty retrieves a single faculty by ID
func GetFaculty(c *gin.Context) {
	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	facultyID := c.Param("facultyID")
	if facultyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "faculty ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	doc, err := f.Collection("faculties").Doc(facultyID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "faculty not found"})
		return
	}

	var faculty models.Faculty
	if err := doc.DataTo(&faculty); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse faculty data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"faculty": faculty,
	})
}

// UpdateFaculty updates an existing faculty
func UpdateFaculty(c *gin.Context) {
	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	facultyID := c.Param("facultyID")
	if facultyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "faculty ID is required"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Dean        string `json:"dean"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build update map
	updates := []firestore.Update{
		{Path: "updatedAt", Value: time.Now()},
	}
	if req.Name != "" {
		updates = append(updates, firestore.Update{Path: "name", Value: req.Name})
	}
	if req.Code != "" {
		updates = append(updates, firestore.Update{Path: "code", Value: req.Code})
	}
	if req.Dean != "" {
		updates = append(updates, firestore.Update{Path: "dean", Value: req.Dean})
	}
	if req.Description != "" {
		updates = append(updates, firestore.Update{Path: "description", Value: req.Description})
	}

	_, err := f.Collection("faculties").Doc(facultyID).Update(ctx, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update faculty"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Faculty updated successfully",
	})
}

// DeleteFaculty deletes a faculty
func DeleteFaculty(c *gin.Context) {
	f, ok := c.MustGet("firestore").(*firestore.Client)
	if !ok || f == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
		return
	}

	facultyID := c.Param("facultyID")
	if facultyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "faculty ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := f.Collection("faculties").Doc(facultyID).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete faculty"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Faculty deleted successfully",
	})
}
