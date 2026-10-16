# Online Course Management System (RESTful Backend)

A robust, performant RESTful backend for an **Online Course Management System** built with **Go**, **Gin Web Framework**, and **SQL** (via GORM with pure-Go SQLite).

---

## 📌 Features & Highlights

- **Domain Entities**:
  - **`Course`**: `id` (Long), `title`, `description`, `durationInHours`, `startDate` (LocalDate format `YYYY-MM-DD`), `instructor`
  - **`Student`**: `id` (Long), `name`, `email` (unique), `enrolledCourses` (Many-to-Many with `Course`)
- **Complete RESTful Endpoints**:
  - Full CRUD for `Course`
  - Full CRUD for `Student`
  - Enroll a student into a course (via path or request body)
  - Unenroll a student from a course
  - Query all courses a student is enrolled in
  - Query all students enrolled in a course
- **Data Seeding**:
  - **CommandLineRunner-style Seeder**: Automatically seeds sample courses, students, and enrollments on initial startup if the database is empty (or on demand).
  - **Pure SQL Seeder**: [`seed/seed.sql`](seed/seed.sql) with DDL and DML scripts.
- **Data Integrity**: Enforces unique constraints, email formatting, and relational cascades.
- **Zero External Server Dependency**: Uses modern pure-Go SQLite (`github.com/glebarez/sqlite`) — runs on any OS without requiring CGO or native compilers.

---

## 🏗️ Architecture & Project Structure

```
course-management-system/
├── config/
│   └── database.go              # Database connection, PRAGMA foreign keys & AutoMigrate
├── controllers/
│   ├── course_controller.go     # Course CRUD and query students in course
│   ├── student_controller.go    # Student CRUD and query student's courses
│   ├── enrollment_controller.go # Many-to-many enrollment and unenrollment
│   └── controllers_test.go      # Integration and unit tests
├── models/
│   ├── course.go                # Course model and CourseInput DTO
│   └── student.go               # Student model, StudentInput & EnrollmentInput DTO
├── routes/
│   └── routes.go                # Route registration and middleware
├── seed/
│   ├── seeder.go                # Programmatic data seeder (CommandLineRunner style)
│   └── seed.sql                 # SQL schema & DML seed script
├── go.mod
├── go.sum
├── main.go                      # Application entrypoint & CLI flags
└── README.md                    # Setup and API documentation
```

---

## 🚀 Getting Started

### Prerequisites
- [Go](https://go.dev/dl/) (version 1.20+ installed)

### 1. Run the Server
From the project root:

```bash
cd course-management-system
go run main.go
```

The server initializes the SQLite database (`course_management.db`), runs migrations, seeds initial sample data automatically, and starts listening on port `8080`.

#### CLI Flags:
- `-port`: Custom port (default: `8080`), e.g., `go run main.go -port 9000`
- `-db`: Custom database file (default: `course_management.db`)
- `-seed`: Force re-seed (default: `false`)

### 2. Run Automated Tests
Execute the comprehensive test suite verifying all CRUD endpoints, uniqueness checks, and many-to-many enrollments:

```bash
go test -v ./...
```

---

## 🗄️ Database Seeding

### Method A: Built-in Startup Seeder (CommandLineRunner Style)
The application automatically checks if the database is populated on startup:
```go
seed.SeedDatabase()
```
Pre-populates:
- 4 Courses (Golang, Database Systems, Cloud-Native, Data Structures)
- 4 Students (Alice, Bob, Charlie, Diana)
- Many-to-many enrollments linking students across courses.

### Method B: SQL Script Execution
You can directly execute the provided SQL script into any SQLite or SQL database client:
```bash
sqlite3 course_management.db < seed/seed.sql
```

---

## 📖 API Documentation & Endpoints

Base URL: `http://localhost:8080`

### 1. Course Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/courses` | Create a new course |
| `GET` | `/api/courses` | Retrieve all courses |
| `GET` | `/api/courses/:id` | Retrieve course details by ID |
| `PUT` | `/api/courses/:id` | Update an existing course |
| `DELETE` | `/api/courses/:id` | Delete a course |
| `GET` | `/api/courses/:id/students` | Get all students enrolled in a course |

#### Create Course (`POST /api/courses`)
**Request Body:**
```json
{
  "title": "Full-Stack Go Development",
  "description": "Building modern full-stack web applications with Go and React.",
  "durationInHours": 45,
  "startDate": "2026-11-20",
  "instructor": "John Doe"
}
```
**Response (201 Created):**
```json
{
  "message": "Course created successfully",
  "data": {
    "id": 5,
    "title": "Full-Stack Go Development",
    "description": "Building modern full-stack web applications with Go and React.",
    "durationInHours": 45,
    "startDate": "2026-11-20",
    "instructor": "John Doe",
    "createdAt": "2026-10-15T11:16:30Z",
    "updatedAt": "2026-10-15T11:16:30Z"
  }
}
```

---

### 2. Student Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/students` | Create a new student |
| `GET` | `/api/students` | Retrieve all students |
| `GET` | `/api/students/:id` | Retrieve student details with enrolled courses |
| `PUT` | `/api/students/:id` | Update an existing student |
| `DELETE` | `/api/students/:id` | Delete a student |
| `GET` | `/api/students/:id/courses` | Get all courses a student is enrolled in |

#### Create Student (`POST /api/students`)
**Request Body:**
```json
{
  "name": "Sarah Connor",
  "email": "sarah.connor@example.com"
}
```
**Response (201 Created):**
```json
{
  "message": "Student created successfully",
  "data": {
    "id": 5,
    "name": "Sarah Connor",
    "email": "sarah.connor@example.com",
    "createdAt": "2026-10-15T11:16:30Z",
    "updatedAt": "2026-10-15T11:16:30Z"
  }
}
```

---

### 3. Enrollment Endpoints (Many-to-Many)

#### Enroll a Student in a Course
You can enroll using either RESTful route style:

**Option 1: Path-based**
```http
POST /api/students/:id/courses/:courseId
```
*Example:* `POST /api/students/1/courses/3`

**Option 2: Payload-based**
```http
POST /api/enrollments
```
**Request Body:**
```json
{
  "studentId": 1,
  "courseId": 3
}
```

**Response (200 OK):**
```json
{
  "message": "Enrollment successful",
  "studentId": 1,
  "studentName": "Alice Smith",
  "courseId": 3,
  "courseTitle": "Cloud-Native Architecture & Kubernetes"
}
```

#### Unenroll a Student from a Course
```http
DELETE /api/students/:id/courses/:courseId
```

**Response (200 OK):**
```json
{
  "message": "Student successfully unenrolled from course",
  "studentId": 1,
  "courseId": 3
}
```

---

### 4. Relationship Queries

#### Get all courses a student is enrolled in:
```http
GET /api/students/1/courses
```
**Response (200 OK):**
```json
{
  "studentId": 1,
  "studentName": "Alice Smith",
  "total": 2,
  "data": [
    {
      "id": 1,
      "title": "Introduction to Golang and Microservices",
      "description": "Learn modern backend engineering with Go, Gin, Goroutines, and REST architecture.",
      "durationInHours": 40,
      "startDate": "2026-11-01",
      "instructor": "Dr. Alan Kay"
    },
    {
      "id": 2,
      "title": "Database Systems and SQL Mastery",
      "description": "Comprehensive guide to Relational databases, ACID compliance, normalization, and indexing.",
      "durationInHours": 35,
      "startDate": "2026-11-10",
      "instructor": "Prof. Barbara Liskov"
    }
  ]
}
```

#### Get all students enrolled in a course:
```http
GET /api/courses/1/students
```
**Response (200 OK):**
```json
{
  "courseId": 1,
  "courseTitle": "Introduction to Golang and Microservices",
  "total": 3,
  "data": [
    {
      "id": 1,
      "name": "Alice Smith",
      "email": "alice.smith@university.edu"
    },
    {
      "id": 2,
      "name": "Bob Jones",
      "email": "bob.jones@university.edu"
    },
    {
      "id": 4,
      "name": "Diana Prince",
      "email": "diana.prince@university.edu"
    }
  ]
}
```

---

## 🧪 Sample cURL Commands

```bash
# 1. Fetch all courses
curl http://localhost:8080/api/courses

# 2. Fetch all students
curl http://localhost:8080/api/students

# 3. Create a new course
curl -X POST http://localhost:8080/api/courses \
  -H "Content-Type: application/json" \
  -d '{"title":"Concurrency in Go","description":"Channels & sync primitives","durationInHours":20,"startDate":"2026-11-25","instructor":"Rob Pike"}'

# 4. Enroll Student 3 in Course 1
curl -X POST http://localhost:8080/api/students/3/courses/1

# 5. List courses for Student 3
curl http://localhost:8080/api/students/3/courses

# 6. List students in Course 1
curl http://localhost:8080/api/courses/1/students
```
