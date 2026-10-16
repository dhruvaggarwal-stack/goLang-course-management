package seed

import (
	"log"

	"course-management-system/config"
	"course-management-system/models"
)

// SeedDatabase acts like Spring's CommandLineRunner, populating sample data if tables are empty
func SeedDatabase() {
	var courseCount int64
	config.DB.Model(&models.Course{}).Count(&courseCount)

	if courseCount > 0 {
		log.Println("[Seeder] Database already contains records. Skipping seed.")
		return
	}

	log.Println("[Seeder] Seeding initial data (courses, students, and enrollments)...")

	// 1. Create sample courses
	c1 := models.Course{
		Title:           "Introduction to Golang and Microservices",
		Description:     "Learn modern backend engineering with Go, Gin, Goroutines, and REST architecture.",
		DurationInHours: 40,
		StartDate:       "2026-11-01",
		Instructor:      "Dr. Alan Kay",
	}
	c2 := models.Course{
		Title:           "Database Systems and SQL Mastery",
		Description:     "Comprehensive guide to Relational databases, ACID compliance, normalization, and indexing.",
		DurationInHours: 35,
		StartDate:       "2026-11-10",
		Instructor:      "Prof. Barbara Liskov",
	}
	c3 := models.Course{
		Title:           "Cloud-Native Architecture & Kubernetes",
		Description:     "Master container orchestration, Helm charts, CI/CD pipelines, and cloud deployment.",
		DurationInHours: 50,
		StartDate:       "2026-12-01",
		Instructor:      "Kelsey Hightower",
	}
	c4 := models.Course{
		Title:           "Data Structures and Algorithms in Go",
		Description:     "Deep dive into trees, graphs, dynamic programming, and algorithmic problem-solving.",
		DurationInHours: 45,
		StartDate:       "2026-11-15",
		Instructor:      "Donald Knuth",
	}

	config.DB.Create(&c1)
	config.DB.Create(&c2)
	config.DB.Create(&c3)
	config.DB.Create(&c4)

	// 2. Create sample students
	s1 := models.Student{
		Name:  "Alice Smith",
		Email: "alice.smith@university.edu",
	}
	s2 := models.Student{
		Name:  "Bob Jones",
		Email: "bob.jones@university.edu",
	}
	s3 := models.Student{
		Name:  "Charlie Brown",
		Email: "charlie.brown@university.edu",
	}
	s4 := models.Student{
		Name:  "Diana Prince",
		Email: "diana.prince@university.edu",
	}

	config.DB.Create(&s1)
	config.DB.Create(&s2)
	config.DB.Create(&s3)
	config.DB.Create(&s4)

	// 3. Create enrollments (Many-to-Many)
	// Alice enlists in Golang and Database Systems
	_ = config.DB.Model(&s1).Association("EnrolledCourses").Append(&c1, &c2)

	// Bob enlists in Golang and Cloud-Native
	_ = config.DB.Model(&s2).Association("EnrolledCourses").Append(&c1, &c3)

	// Charlie enlists in Database Systems and Algorithms
	_ = config.DB.Model(&s3).Association("EnrolledCourses").Append(&c2, &c4)

	// Diana enlists in all four courses
	_ = config.DB.Model(&s4).Association("EnrolledCourses").Append(&c1, &c2, &c3, &c4)

	log.Println("[Seeder] Successfully seeded 4 courses, 4 students, and cross-course enrollments.")
}
