package controllers

import (
	"net/http"
	"strconv"
	"time"

	"course-management-system/config"
	"course-management-system/models"

	"github.com/gin-gonic/gin"
)

// validateDate checks if date string adheres to YYYY-MM-DD
func isValidDateFormat(d string) bool {
	_, err := time.Parse("2006-01-02", d)
	return err == nil
}

// CreateCourse godoc
// @Summary Create a new course
// @Route POST /api/courses
func CreateCourse(c *gin.Context) {
	var input models.CourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	if !isValidDateFormat(input.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid startDate format. Expected YYYY-MM-DD",
		})
		return
	}

	course := models.Course{
		Title:           input.Title,
		Description:     input.Description,
		DurationInHours: input.DurationInHours,
		StartDate:       input.StartDate,
		Instructor:      input.Instructor,
	}

	if err := config.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Course created successfully",
		"data":    course,
	})
}

// GetAllCourses godoc
// @Summary Get all courses
// @Route GET /api/courses
func GetAllCourses(c *gin.Context) {
	var courses []models.Course
	// We can choose to preload students or keep it concise
	if err := config.DB.Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(courses),
		"data":  courses,
	})
}

// GetCourseByID godoc
// @Summary Get course by ID
// @Route GET /api/courses/:id
func GetCourseByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var course models.Course
	if err := config.DB.Preload("Students").First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": course,
	})
}

// UpdateCourse godoc
// @Summary Update an existing course
// @Route PUT /api/courses/:id
func UpdateCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var course models.Course
	if err := config.DB.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	var input models.CourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	if !isValidDateFormat(input.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid startDate format. Expected YYYY-MM-DD",
		})
		return
	}

	course.Title = input.Title
	course.Description = input.Description
	course.DurationInHours = input.DurationInHours
	course.StartDate = input.StartDate
	course.Instructor = input.Instructor

	if err := config.DB.Save(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Course updated successfully",
		"data":    course,
	})
}

// DeleteCourse godoc
// @Summary Delete course by ID
// @Route DELETE /api/courses/:id
func DeleteCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var course models.Course
	if err := config.DB.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// Remove associations first from join table, then delete course
	_ = config.DB.Model(&course).Association("Students").Clear()

	if err := config.DB.Delete(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Course deleted successfully",
		"id":      id,
	})
}

// GetStudentsInCourse godoc
// @Summary Get all students enrolled in a course
// @Route GET /api/courses/:id/students
func GetStudentsInCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var course models.Course
	if err := config.DB.Preload("Students").First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courseId":    course.ID,
		"courseTitle": course.Title,
		"total":       len(course.Students),
		"data":        course.Students,
	})
}
