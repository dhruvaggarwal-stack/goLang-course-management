package controllers

import (
	"net/http"
	"strconv"

	"course-management-system/config"
	"course-management-system/models"
	"github.com/gin-gonic/gin"
)

// EnrollStudentByPath godoc
// @Summary Enroll a student in a course via URL path
// @Route POST /api/students/:id/courses/:courseId
func EnrollStudentByPath(c *gin.Context) {
	studentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	enroll(c, uint(studentID), uint(courseID))
}

// EnrollStudentByBody godoc
// @Summary Enroll a student in a course via JSON body
// @Route POST /api/enrollments
func EnrollStudentByBody(c *gin.Context) {
	var input models.EnrollmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	enroll(c, input.StudentID, input.CourseID)
}

// UnenrollStudent godoc
// @Summary Remove student enrollment from a course
// @Route DELETE /api/students/:id/courses/:courseId
func UnenrollStudent(c *gin.Context) {
	studentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var student models.Student
	if err := config.DB.First(&student, studentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	var course models.Course
	if err := config.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// Remove relationship
	if err := config.DB.Model(&student).Association("EnrolledCourses").Delete(&course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unenroll student", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Student successfully unenrolled from course",
		"studentId": studentID,
		"courseId":  courseID,
	})
}

// Helper function to handle enrollment logic
func enroll(c *gin.Context, studentID uint, courseID uint) {
	var student models.Student
	if err := config.DB.First(&student, studentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	var course models.Course
	if err := config.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// Check if already enrolled
	var count int64
	config.DB.Table("student_courses").
		Where("student_id = ? AND course_id = ?", studentID, courseID).
		Count(&count)

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":     "Student is already enrolled in this course",
			"studentId": studentID,
			"courseId":  courseID,
		})
		return
	}

	// Add association
	if err := config.DB.Model(&student).Association("EnrolledCourses").Append(&course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enroll student", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Enrollment successful",
		"studentId":   student.ID,
		"studentName": student.Name,
		"courseId":    course.ID,
		"courseTitle": course.Title,
	})
}
