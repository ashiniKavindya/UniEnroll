package models

import "time"

// EnrollmentStatus represents the status of module enrollment
const (
	EnrollmentStatusPending   = "pending"
	EnrollmentStatusEnrolled  = "enrolled"
	EnrollmentStatusCompleted = "completed"
	EnrollmentStatusDropped   = "dropped"
	EnrollmentStatusFailed    = "failed"
)

// Enrollment represents a student's enrollment in a module
type Enrollment struct {
	EnrollmentID string    `firestore:"enrollmentID" json:"enrollmentID"`
	StudentID    string    `firestore:"studentID" json:"studentID"`
	ModuleID     string    `firestore:"moduleID" json:"moduleID"`
	Status       string    `firestore:"status" json:"status"` // pending, enrolled, completed, dropped, failed
	Grade        string    `firestore:"grade,omitempty" json:"grade,omitempty"`
	Score        float64   `firestore:"score,omitempty" json:"score,omitempty"`
	EnrolledAt   time.Time `firestore:"enrolledAt" json:"enrolledAt"`
	CompletedAt  time.Time `firestore:"completedAt,omitempty" json:"completedAt,omitempty"`
	CreatedAt    time.Time `firestore:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time `firestore:"updatedAt" json:"updatedAt"`
}
