package models

// User represents a user stored in Firestore
type User struct {
	Name         string `firestore:"name,omitempty" json:"name"`
	Email        string `firestore:"email,omitempty" json:"email"`
	PasswordHash string `firestore:"passwordHash,omitempty" json:"-"`
	Role         string `firestore:"role,omitempty" json:"role"`
}
