package models

import "time"

// Lecturer represents a lecturer in the system
type Lecturer struct {
	UserID         string    `firestore:"userID" json:"userID"`
	Department     string    `firestore:"department,omitempty" json:"department"`
	Specialization string    `firestore:"specialization,omitempty" json:"specialization"`
	OfficeLocation string    `firestore:"officeLocation,omitempty" json:"officeLocation"`
	PhoneNumber    string    `firestore:"phoneNumber,omitempty" json:"phoneNumber"`
	Qualifications []string  `firestore:"qualifications,omitempty" json:"qualifications"`
	CreatedAt      time.Time `firestore:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time `firestore:"updatedAt" json:"updatedAt"`
}
