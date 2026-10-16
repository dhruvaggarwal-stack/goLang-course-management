package main

import (
	"flag"
	"fmt"
	"log"

	"course-management-system/config"
	"course-management-system/routes"
	"course-management-system/seed"
)

func main() {
	port := flag.String("port", "8080", "Server port to listen on")
	dbPath := flag.String("db", "course_management.db", "SQLite database file path")
	forceSeed := flag.Bool("seed", false, "Force seed sample data into database")
	flag.Parse()

	log.Println("Starting Online Course Management System Backend...")

	// 1. Initialize SQLite Database & Auto-Migrate
	config.ConnectDatabase(*dbPath)

	// 2. Seed Sample Data (CommandLineRunner style)
	if *forceSeed {
		log.Println("Force seed flag detected.")
	}
	seed.SeedDatabase()

	// 3. Initialize Gin Router
	r := routes.SetupRouter()

	address := fmt.Sprintf(":%s", *port)
	log.Printf("Server is running and listening at http://localhost%s\n", address)
	log.Println("Available endpoints:")
	log.Println(" - GET    /api/courses")
	log.Println(" - POST   /api/courses")
	log.Println(" - GET    /api/courses/:id")
	log.Println(" - PUT    /api/courses/:id")
	log.Println(" - DELETE /api/courses/:id")
	log.Println(" - GET    /api/courses/:id/students")
	log.Println(" - GET    /api/students")
	log.Println(" - POST   /api/students")
	log.Println(" - GET    /api/students/:id")
	log.Println(" - PUT    /api/students/:id")
	log.Println(" - DELETE /api/students/:id")
	log.Println(" - GET    /api/students/:id/courses")
	log.Println(" - POST   /api/students/:id/courses/:courseId")
	log.Println(" - DELETE /api/students/:id/courses/:courseId")
	log.Println(" - POST   /api/enrollments")

	if err := r.Run(address); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
