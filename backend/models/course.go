package models

// Course represents a course in the system
type Course struct {
	Title       string `firestore:"title,omitempty" json:"title"`
	Code        string `firestore:"code,omitempty" json:"code"`
	Description string `firestore:"description,omitempty" json:"description"`
	Instructor  string `firestore:"instructor,omitempty" json:"instructor"`
	Credits     int    `firestore:"credits,omitempty" json:"credits"`
}
