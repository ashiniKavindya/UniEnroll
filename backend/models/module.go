package models

// Module represents a module/subject within a course
type Module struct {
	Title       string `firestore:"title,omitempty" json:"title"`
	CourseID    string `firestore:"courseId,omitempty" json:"courseId"`
	Description string `firestore:"description,omitempty" json:"description"`
	Credits     int    `firestore:"credits,omitempty" json:"credits"`
}
