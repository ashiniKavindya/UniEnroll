package models

import "time"

// Student represents a student in the system
type Student struct {
	UserID         string    `firestore:"userID" json:"userID"`
	StudentNumber  string    `firestore:"studentNumber,omitempty" json:"studentNumber"`
	DepartmentID   string    `firestore:"departmentID,omitempty" json:"departmentID"`
	YearOfStudy    int       `firestore:"yearOfStudy,omitempty" json:"yearOfStudy"`
	EnrollmentYear int       `firestore:"enrollmentYear,omitempty" json:"enrollmentYear"`
	CreatedAt      time.Time `firestore:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time `firestore:"updatedAt" json:"updatedAt"`
}
