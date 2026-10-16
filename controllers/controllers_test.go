package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"course-management-system/config"
	"course-management-system/models"
	"course-management-system/routes"
	"course-management-system/seed"
	"github.com/gin-gonic/gin"
)

var testRouter *gin.Engine

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	// Use temporary SQLite database for testing
	testDBPath := "test_course_mgmt.db"
	_ = os.Remove(testDBPath)

	config.ConnectDatabase(testDBPath)
	seed.SeedDatabase()

	testRouter = routes.SetupRouter()

	exitCode := m.Run()

	_ = os.Remove(testDBPath)
	os.Exit(exitCode)
}

func TestGetAllCourses(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/api/courses", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var res map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	count := int(res["count"].(float64))
	if count < 4 {
		t.Fatalf("Expected at least 4 seeded courses, got %d", count)
	}
}

func TestCreateCourse(t *testing.T) {
	body := models.CourseInput{
		Title:           "Distributed Systems in Go",
		Description:     "Raft, Paxos, Consensus, and Event-Driven architecture",
		DurationInHours: 60,
		StartDate:       "2026-12-15",
		Instructor:      "Leslie Lamport",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/courses", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 Created, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestGetAllStudents(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/api/students", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}

func TestCreateStudentDuplicateEmail(t *testing.T) {
	// alice.smith@university.edu is already seeded
	body := models.StudentInput{
		Name:  "Alice Clone",
		Email: "alice.smith@university.edu",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/students", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("Expected status 409 Conflict for duplicate email, got %d", w.Code)
	}
}

func TestGetStudentCourses(t *testing.T) {
	// Student 1 (Alice) is enrolled in courses
	req, _ := http.NewRequest(http.MethodGet, "/api/students/1/courses", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var res map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	total := int(res["total"].(float64))
	if total == 0 {
		t.Fatalf("Expected student 1 to have enrolled courses, got 0")
	}
}

func TestGetStudentsInCourse(t *testing.T) {
	// Course 1 has students enrolled
	req, _ := http.NewRequest(http.MethodGet, "/api/courses/1/students", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var res map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	total := int(res["total"].(float64))
	if total == 0 {
		t.Fatalf("Expected course 1 to have students, got 0")
	}
}

func TestEnrollAndUnenrollStudent(t *testing.T) {
	// Student 2 (Bob) enlisting in Course 2 (Database Systems)
	req, _ := http.NewRequest(http.MethodPost, "/api/students/2/courses/2", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK for enrollment, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Repeat enrollment should return 409 Conflict
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodPost, "/api/students/2/courses/2", nil)
	testRouter.ServeHTTP(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Fatalf("Expected status 409 Conflict for duplicate enrollment, got %d", w2.Code)
	}

	// Unenroll
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest(http.MethodDelete, "/api/students/2/courses/2", nil)
	testRouter.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK for unenrollment, got %d", w3.Code)
	}
}
