package models

import "time"

// Department represents an academic department
type Department struct {
	DepartmentID string    `firestore:"departmentID" json:"departmentID"`
	Name         string    `firestore:"name" json:"name"`
	Code         string    `firestore:"code" json:"code"`
	FacultyName  string    `firestore:"facultyName,omitempty" json:"facultyName"`
	HOD          string    `firestore:"hod,omitempty" json:"hod"` // Head of Department user ID
	Description  string    `firestore:"description,omitempty" json:"description"`
	CreatedAt    time.Time `firestore:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time `firestore:"updatedAt" json:"updatedAt"`
}
