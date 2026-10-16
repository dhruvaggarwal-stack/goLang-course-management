package models

import (
	"time"
)

// Student represents a student entity
type Student struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name" binding:"required"`
	Email     string    `gorm:"size:255;not null;uniqueIndex" json:"email" binding:"required,email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Many-to-Many relationship with Course
	EnrolledCourses []Course `gorm:"many2many:student_courses;constraint:OnDelete:CASCADE;" json:"enrolledCourses,omitempty"`
}

// StudentInput is used for creating or updating a student
type StudentInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// EnrollmentInput is used when enrolling via JSON body
type EnrollmentInput struct {
	StudentID uint `json:"studentId" binding:"required"`
	CourseID  uint `json:"courseId" binding:"required"`
}
