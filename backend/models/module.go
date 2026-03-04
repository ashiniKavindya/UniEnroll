package models

import "time"

// ModuleType represents the type of module
const (
	ModuleTypeCompulsory = "compulsory"
	ModuleTypeElective   = "elective"
	ModuleTypeFaculty    = "faculty" // Module from another department
)

// Module represents a module/subject within a course
type Module struct {
	ModuleID      string    `firestore:"moduleID" json:"moduleID"`
	Title         string    `firestore:"title,omitempty" json:"title"`
	Code          string    `firestore:"code,omitempty" json:"code"`
	DepartmentID  string    `firestore:"departmentID,omitempty" json:"departmentID"` // Owning department
	CourseID      string    `firestore:"courseID,omitempty" json:"courseID"`
	LecturerID    string    `firestore:"lecturerID,omitempty" json:"lecturerID"`
	Type          string    `firestore:"type,omitempty" json:"type"` // compulsory, elective, faculty
	Description   string    `firestore:"description,omitempty" json:"description"`
	Credits       int       `firestore:"credits,omitempty" json:"credits"`
	Semester      int       `firestore:"semester,omitempty" json:"semester"`
	YearOfStudy   int       `firestore:"yearOfStudy,omitempty" json:"yearOfStudy"`     // Which year can take this
	Prerequisites []string  `firestore:"prerequisites,omitempty" json:"prerequisites"` // Module IDs required before taking this
	CreatedAt     time.Time `firestore:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time `firestore:"updatedAt" json:"updatedAt"`
}
