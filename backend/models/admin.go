package models

import "time"

// Admin represents an admin user in the system
type Admin struct {
	UserID      string    `firestore:"userID" json:"userID"`
	Permissions []string  `firestore:"permissions,omitempty" json:"permissions"`
	Department  string    `firestore:"department,omitempty" json:"department"`
	PhoneNumber string    `firestore:"phoneNumber,omitempty" json:"phoneNumber"`
	CreatedAt   time.Time `firestore:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `firestore:"updatedAt" json:"updatedAt"`
}
