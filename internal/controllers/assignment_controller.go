package controllers

import (
	"context"

	m "github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type AssignmentController struct {
	DB *pgxpool.Pool
	lg *logrus.Logger

	userc *UserController
}

func NewAssignmentController(db *pgxpool.Pool, lg *logrus.Logger, userc *UserController) *AssignmentController {
	return &AssignmentController{
		DB:    db,
		lg:    lg,
		userc: userc,
	}
}

func (assigh *AssignmentController) CreateAssignment(ctx context.Context, assignment *m.Assignment, studentIDs []int, userID int) error {

	user, err := assigh.userc.GetUserByID(ctx, userID)
	if err != nil {
		assigh.lg.Printf("user controller - CreateAssignment - get user by id - %v", err)
		return err
	}

	if user.Role != "teacher" {
		assigh.lg.Println("assignment controller - CreateAssignment - no permission for user")
		return nil
	}
	assigh.lg.Debugln("Creating assignment at controller level")
	tx, err := assigh.DB.Begin(ctx)
	if err != nil {
		assigh.lg.Errorf("user controller - CreateAssignment - db begin transaction - %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	var assignmentID int
	// Insert the assignment and get the ID
	err = tx.QueryRow(ctx, "INSERT INTO assignments (name, topic_id, weight, teacher_id) VALUES ($1, $2, $3, $4) RETURNING id",
		assignment.Name, assignment.TopicID, assignment.Weight, assignment.TeacherID).Scan(&assignmentID)
	if err != nil {
		assigh.lg.Errorf("user controller - CreateAssignment - db exec - %v", err)
		return err
	}

	// Link the assignment to students
	for _, studentID := range studentIDs {
		_, err = tx.Exec(ctx, "INSERT INTO assignment_students (assignment_id, student_id) VALUES ($1, $2)",
			assignmentID, studentID)
		if err != nil {
			assigh.lg.Errorf("user controller - CreateAssignment - link to students - %v", err)
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		assigh.lg.Errorf("user controller - CreateAssignment - commit transaction - %v", err)
		return err
	}

	return nil
}

func (assigh *AssignmentController) UpdateAssignment(ctx context.Context, assignment *m.Assignment) *m.Assignment {
	assigh.lg.Debugln("Updating assignment at controller level")
	tx, err := assigh.DB.Begin(ctx)
	if err != nil {
		assigh.lg.Errorf("assignment controller - UpdateAssignment - db begin transaction - %v", err)
		return nil // Return nil on error
	}
	defer tx.Rollback(ctx)

	// Update the assignment based on the provided fields
	result, err := tx.Exec(ctx, `
        UPDATE assignments 
        SET name = $1, topic_id = $2, weight = $3, teacher_id = $4
        WHERE id = $5
    `, assignment.Name, assignment.TopicID, assignment.Weight, assignment.TeacherID, assignment.ID)

	if err != nil {
		assigh.lg.Errorf("assignment controller - UpdateAssignment - db exec - %v", err)
		return nil
	}

	// Handle cases where the assignment is not found
	if result.RowsAffected() == 0 {
		assigh.lg.Warnf("assignment controller - UpdateAssignment - assignment not found (ID: %d)", assignment.ID)

		// You can decide to either:
		// 1. Return nil to signal the record was not found
		return nil
		// 2. Create a new assignment if you deem it appropriate
	}

	if err = tx.Commit(ctx); err != nil {
		assigh.lg.Errorf("assignment controller - UpdateAssignment - commit transaction - %v", err)
		return nil
	}

	// Return the updated assignment (you'll likely need to fetch it from the database again)
	return assignment
}
