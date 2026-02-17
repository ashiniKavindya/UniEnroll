package models

import "time"

// Enrollment represents a student's enrollment in a course
type Enrollment struct {
	UserID     string    `firestore:"userId,omitempty" json:"userId"`
	CourseID   string    `firestore:"courseId,omitempty" json:"courseId"`
	EnrolledAt time.Time `firestore:"enrolledAt,omitempty" json:"enrolledAt"`
	Status     string    `firestore:"status,omitempty" json:"status"` // "active", "completed", "dropped"
}
