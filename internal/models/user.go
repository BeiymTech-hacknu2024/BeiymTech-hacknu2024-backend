package models

import (
	"time"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Role     string
}

type UserActivity struct {
	ID     int
	UserID int
	Route  string
	Time   time.Time
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type StudentAssignment struct {
	ID           int
	AssignmentID int
	StudentID    int
	Score        *int // Using a pointer to int to allow null (not completed assignments)
}

type StudentPerformanceBySubject struct {
	SubjectID    int
	StudentID    int
	OverallScore *int // Using a pointer to int to allow null values
}
