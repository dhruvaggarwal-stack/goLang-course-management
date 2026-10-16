-- ========================================================
-- Online Course Management System - SQL Seed Script
-- ========================================================

-- Enable Foreign Keys
PRAGMA foreign_keys = ON;

-- 1. Create Courses Table
CREATE TABLE IF NOT EXISTS courses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration_in_hours INTEGER NOT NULL,
    start_date VARCHAR(10) NOT NULL,
    instructor VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2. Create Students Table
CREATE TABLE IF NOT EXISTS students (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 3. Create Student-Courses Join Table (Many-to-Many)
CREATE TABLE IF NOT EXISTS student_courses (
    student_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    PRIMARY KEY (student_id, course_id),
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- 4. Sample Courses Data
INSERT OR IGNORE INTO courses (id, title, description, duration_in_hours, start_date, instructor) VALUES
(1, 'Introduction to Golang and Microservices', 'Learn modern backend engineering with Go, Gin, Goroutines, and REST architecture.', 40, '2026-11-01', 'Dr. Alan Kay'),
(2, 'Database Systems and SQL Mastery', 'Comprehensive guide to Relational databases, ACID compliance, normalization, and indexing.', 35, '2026-11-10', 'Prof. Barbara Liskov'),
(3, 'Cloud-Native Architecture & Kubernetes', 'Master container orchestration, Helm charts, CI/CD pipelines, and cloud deployment.', 50, '2026-12-01', 'Kelsey Hightower'),
(4, 'Data Structures and Algorithms in Go', 'Deep dive into trees, graphs, dynamic programming, and algorithmic problem-solving.', 45, '2026-11-15', 'Donald Knuth');

-- 5. Sample Students Data
INSERT OR IGNORE INTO students (id, name, email) VALUES
(1, 'Alice Smith', 'alice.smith@university.edu'),
(2, 'Bob Jones', 'bob.jones@university.edu'),
(3, 'Charlie Brown', 'charlie.brown@university.edu'),
(4, 'Diana Prince', 'diana.prince@university.edu');

-- 6. Sample Enrollments Data (Many-to-Many Mappings)
INSERT OR IGNORE INTO student_courses (student_id, course_id) VALUES
(1, 1), -- Alice in Golang
(1, 2), -- Alice in Database Systems
(2, 1), -- Bob in Golang
(2, 3), -- Bob in Cloud-Native
(3, 2), -- Charlie in Database Systems
(3, 4), -- Charlie in Algorithms
(4, 1), -- Diana in Golang
(4, 2), -- Diana in Database Systems
(4, 3), -- Diana in Cloud-Native
(4, 4); -- Diana in Algorithms
