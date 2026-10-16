package routes

import (
	"course-management-system/controllers"
	"github.com/gin-gonic/gin"
)

// SetupRouter sets up and registers all the endpoints
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS and recovery middleware
	r.Use(gin.Recovery())

	// Welcome & Health endpoints
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the Online Course Management System API",
			"status":  "running",
			"docs":    "/api/docs",
		})
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	api := r.Group("/api")
	{
		// Course Endpoints
		courses := api.Group("/courses")
		{
			courses.POST("", controllers.CreateCourse)
			courses.GET("", controllers.GetAllCourses)
			courses.GET("/:id", controllers.GetCourseByID)
			courses.PUT("/:id", controllers.UpdateCourse)
			courses.DELETE("/:id", controllers.DeleteCourse)
			courses.GET("/:id/students", controllers.GetStudentsInCourse)
		}

		// Student Endpoints
		students := api.Group("/students")
		{
			students.POST("", controllers.CreateStudent)
			students.GET("", controllers.GetAllStudents)
			students.GET("/:id", controllers.GetStudentByID)
			students.PUT("/:id", controllers.UpdateStudent)
			students.DELETE("/:id", controllers.DeleteStudent)
			students.GET("/:id/courses", controllers.GetStudentCourses)

			// Enrollment route under student resource
			students.POST("/:id/courses/:courseId", controllers.EnrollStudentByPath)
			students.DELETE("/:id/courses/:courseId", controllers.UnenrollStudent)
		}

		// Enrollment Endpoints (Alternative POST /api/enrollments)
		api.POST("/enrollments", controllers.EnrollStudentByBody)
	}

	return r
}
