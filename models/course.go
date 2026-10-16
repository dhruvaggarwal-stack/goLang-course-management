package models

import (
	"time"
)

// Course represents a course entity
type Course struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string    `gorm:"size:255;not null" json:"title" binding:"required"`
	Description     string    `gorm:"type:text" json:"description"`
	DurationInHours int       `gorm:"not null" json:"durationInHours" binding:"required,min=1"`
	StartDate       string    `gorm:"type:varchar(10);not null" json:"startDate" binding:"required"` // Format: YYYY-MM-DD
	Instructor      string    `gorm:"size:255;not null" json:"instructor" binding:"required"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	// Many-to-Many relationship with Student
	Students []Student `gorm:"many2many:student_courses;constraint:OnDelete:CASCADE;" json:"students,omitempty"`
}

// CourseInput is used for creating or updating courses
type CourseInput struct {
	Title           string `json:"title" binding:"required"`
	Description     string `json:"description"`
	DurationInHours int    `json:"durationInHours" binding:"required,min=1"`
	StartDate       string `json:"startDate" binding:"required"` // Format: YYYY-MM-DD
	Instructor      string `json:"instructor" binding:"required"`
}
