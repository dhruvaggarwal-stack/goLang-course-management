package controllers

import (
	"net/http"
	"strconv"

	"course-management-system/config"
	"course-management-system/models"
	"github.com/gin-gonic/gin"
)

// CreateStudent godoc
// @Summary Create a new student
// @Route POST /api/students
func CreateStudent(c *gin.Context) {
	var input models.StudentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Check for duplicate email
	var existing models.Student
	if err := config.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "A student with this email already exists"})
		return
	}

	student := models.Student{
		Name:  input.Name,
		Email: input.Email,
	}

	if err := config.DB.Create(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Student created successfully",
		"data":    student,
	})
}

// GetAllStudents godoc
// @Summary Get all students
// @Route GET /api/students
func GetAllStudents(c *gin.Context) {
	var students []models.Student
	if err := config.DB.Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(students),
		"data":  students,
	})
}

// GetStudentByID godoc
// @Summary Get student by ID
// @Route GET /api/students/:id
func GetStudentByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var student models.Student
	if err := config.DB.Preload("EnrolledCourses").First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": student,
	})
}

// UpdateStudent godoc
// @Summary Update an existing student
// @Route PUT /api/students/:id
func UpdateStudent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var student models.Student
	if err := config.DB.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	var input models.StudentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Verify email uniqueness if email changed
	if input.Email != student.Email {
		var duplicate models.Student
		if err := config.DB.Where("email = ? AND id != ?", input.Email, id).First(&duplicate).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email is already taken by another student"})
			return
		}
	}

	student.Name = input.Name
	student.Email = input.Email

	if err := config.DB.Save(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Student updated successfully",
		"data":    student,
	})
}

// DeleteStudent godoc
// @Summary Delete student by ID
// @Route DELETE /api/students/:id
func DeleteStudent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var student models.Student
	if err := config.DB.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// Clear many-to-many associations before deletion
	_ = config.DB.Model(&student).Association("EnrolledCourses").Clear()

	if err := config.DB.Delete(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Student deleted successfully",
		"id":      id,
	})
}

// GetStudentCourses godoc
// @Summary Get all courses a student is enrolled in
// @Route GET /api/students/:id/courses
func GetStudentCourses(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var student models.Student
	if err := config.DB.Preload("EnrolledCourses").First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"studentId":   student.ID,
		"studentName": student.Name,
		"total":       len(student.EnrolledCourses),
		"data":        student.EnrolledCourses,
	})
}
