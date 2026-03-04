package models

import "time"

// Faculty represents an academic faculty (e.g., Faculty of Science, Faculty of Engineering)
type Faculty struct {
	FacultyID   string    `firestore:"facultyID" json:"facultyID"`
	Name        string    `firestore:"name" json:"name"`
	Code        string    `firestore:"code" json:"code"`
	Dean        string    `firestore:"dean,omitempty" json:"dean"` // Dean user ID
	Description string    `firestore:"description,omitempty" json:"description"`
	CreatedAt   time.Time `firestore:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `firestore:"updatedAt" json:"updatedAt"`
}
