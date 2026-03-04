package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"golang.org/x/crypto/bcrypt"
)

// InitializeCollections ensures all required collections exist with proper indexes
func InitializeCollections(ctx context.Context, client *firestore.Client) error {
	log.Println("Initializing Firestore collections...")

	// Collections to create (Firestore creates collections automatically on first write)
	collections := []string{
		"users",
		"students",
		"lecturers",
		"faculties",
		"departments",
		"courses",
		"modules",
		"enrollments",
		"admins",
	}

	for _, col := range collections {
		log.Printf("Ensuring collection '%s' exists...", col)
		// Collections are created on first write, but we can at least verify by checking root
		if err := verifyOrCreateCollection(ctx, client, col); err != nil {
			return fmt.Errorf("failed to initialize collection %s: %v", col, err)
		}
	}

	log.Println("✅ All collections initialized successfully!")
	return nil
}

// verifyOrCreateCollection checks if a collection exists, creates it if needed
func verifyOrCreateCollection(ctx context.Context, client *firestore.Client, name string) error {
	col := client.Collection(name)
	
	// Try to list documents (this will work even if collection is empty)
	docs, err := col.Limit(1).Documents(ctx).GetAll()
	if err != nil {
		return err
	}
	
	// If no docs exist, create a temporary document to ensure collection exists
	if len(docs) == 0 {
		_, err := col.Doc("_init").Set(ctx, map[string]interface{}{
			"initialized": true,
			"timestamp":   time.Now(),
		})
		if err != nil {
			return err
		}
		
		// Delete the temporary document
		_, err = col.Doc("_init").Delete(ctx)
		if err != nil {
			log.Printf("Warning: failed to delete temporary init doc for %s: %v", name, err)
		}
	}
	
	return nil
}

// SeedInitialData populates Firestore with sample data
func SeedInitialData(ctx context.Context, client *firestore.Client) error {
	log.Println("Seeding initial data...")

	// Seed users
	if err := seedUsers(ctx, client); err != nil {
		return fmt.Errorf("failed to seed users: %v", err)
	}

	// Seed faculties
	if err := seedFaculties(ctx, client); err != nil {
		return fmt.Errorf("failed to seed faculties: %v", err)
	}

	// Seed departments
	if err := seedDepartments(ctx, client); err != nil {
		return fmt.Errorf("failed to seed departments: %v", err)
	}

	// Seed courses
	if err := seedCourses(ctx, client); err != nil {
		return fmt.Errorf("failed to seed courses: %v", err)
	}

	// Seed modules
	if err := seedModules(ctx, client); err != nil {
		return fmt.Errorf("failed to seed modules: %v", err)
	}

	log.Println("✅ Initial data seeded successfully!")
	return nil
}

func seedUsers(ctx context.Context, client *firestore.Client) error {
	adminHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), 10)
	if err != nil {
		return err
	}
	lecturerHash, err := bcrypt.GenerateFromPassword([]byte("lecturer123"), 10)
	if err != nil {
		return err
	}
	studentHash, err := bcrypt.GenerateFromPassword([]byte("student123"), 10)
	if err != nil {
		return err
	}

	users := []interface{}{
		map[string]interface{}{
			"name":                  "Admin User",
			"email":                 "admin@uni.edu",
			"passwordHash":          string(adminHash),
			"role":                  "admin",
			"forcePasswordChange":   false,
		},
		map[string]interface{}{
			"name":                  "John Lecturer",
			"email":                 "john@uni.edu",
			"passwordHash":          string(lecturerHash),
			"role":                  "lecturer",
			"forcePasswordChange":   false,
		},
		map[string]interface{}{
			"name":                  "Jane Student",
			"email":                 "jane@uni.edu",
			"passwordHash":          string(studentHash),
			"role":                  "student",
			"forcePasswordChange":   false,
		},
	}

	col := client.Collection("users")
	for _, user := range users {
		_, err = col.Doc(user.(map[string]interface{})["email"].(string)).Set(ctx, user)
		if err != nil {
			return err
		}
	}
	log.Printf("✓ Seeded %d users", len(users))
	return nil
}

func seedFaculties(ctx context.Context, client *firestore.Client) error {
	faculties := []map[string]interface{}{
		{
			"name":        "Faculty of Science",
			"code":        "FS",
			"description": "Science and Engineering",
			"createdAt":   time.Now(),
			"updatedAt":   time.Now(),
		},
		{
			"name":        "Faculty of Arts",
			"code":        "FA",
			"description": "Humanities and Social Sciences",
			"createdAt":   time.Now(),
			"updatedAt":   time.Now(),
		},
	}

	col := client.Collection("faculties")
	for i, faculty := range faculties {
		_, err := col.Doc(fmt.Sprintf("faculty_%d", i+1)).Set(ctx, faculty)
		if err != nil {
			return err
		}
	}
	log.Printf("✓ Seeded %d faculties", len(faculties))
	return nil
}

func seedDepartments(ctx context.Context, client *firestore.Client) error {
	departments := []map[string]interface{}{
		{
			"departmentID": "dept_001",
			"name":         "Computer Science",
			"code":         "CS",
			"facultyID":    "faculty_1",
			"description":  "Computer Science Department",
			"createdAt":    time.Now(),
			"updatedAt":    time.Now(),
		},
		{
			"departmentID": "dept_002",
			"name":         "Mathematics",
			"code":         "MATH",
			"facultyID":    "faculty_1",
			"description":  "Mathematics Department",
			"createdAt":    time.Now(),
			"updatedAt":    time.Now(),
		},
	}

	col := client.Collection("departments")
	for _, dept := range departments {
		_, err := col.Doc(dept["departmentID"].(string)).Set(ctx, dept)
		if err != nil {
			return err
		}
	}
	log.Printf("✓ Seeded %d departments", len(departments))
	return nil
}

func seedCourses(ctx context.Context, client *firestore.Client) error {
	courses := []map[string]interface{}{
		{
			"title":       "Introduction to Programming",
			"code":        "CS101",
			"description": "Learn the basics of programming",
			"instructor":  "john@uni.edu",
			"credits":     3,
		},
		{
			"title":       "Data Structures",
			"code":        "CS201",
			"description": "Advanced data structures and algorithms",
			"instructor":  "john@uni.edu",
			"credits":     4,
		},
	}

	col := client.Collection("courses")
	for i, course := range courses {
		_, err := col.Doc(fmt.Sprintf("course_%d", i+1)).Set(ctx, course)
		if err != nil {
			return err
		}
	}
	log.Printf("✓ Seeded %d courses", len(courses))
	return nil
}

func seedModules(ctx context.Context, client *firestore.Client) error {
	modules := []map[string]interface{}{
		{
			"moduleID":    "mod_001",
			"title":       "Python Basics",
			"code":        "CS101-MOD1",
			"courseID":    "course_1",
			"description": "Python programming fundamentals",
			"credits":     3,
			"createdAt":   time.Now(),
			"updatedAt":   time.Now(),
		},
		{
			"moduleID":    "mod_002",
			"title":       "Web Development",
			"code":        "CS101-MOD2",
			"courseID":    "course_1",
			"description": "Web development with modern frameworks",
			"credits":     3,
			"createdAt":   time.Now(),
			"updatedAt":   time.Now(),
		},
	}

	col := client.Collection("modules")
	for _, module := range modules {
		_, err := col.Doc(module["moduleID"].(string)).Set(ctx, module)
		if err != nil {
			return err
		}
	}
	log.Printf("✓ Seeded %d modules", len(modules))
	return nil
}

// ClearAllData deletes all documents in all collections (useful for testing)
func ClearAllData(ctx context.Context, client *firestore.Client) error {
	collections := []string{
		"users",
		"students",
		"lecturers",
		"faculties",
		"departments",
		"courses",
		"modules",
		"enrollments",
		"admins",
	}

	for _, colName := range collections {
		col := client.Collection(colName)
		docs, err := col.Documents(ctx).GetAll()
		if err != nil {
			return fmt.Errorf("failed to get documents from %s: %v", colName, err)
		}

		for _, doc := range docs {
			if _, err := doc.Ref.Delete(ctx); err != nil {
				return fmt.Errorf("failed to delete document %s from %s: %v", doc.Ref.ID, colName, err)
			}
		}
		log.Printf("✓ Cleared collection '%s'", colName)
	}

	return nil
}
