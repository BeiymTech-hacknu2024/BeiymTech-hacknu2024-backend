package controllers

import (
	"context"
	"fmt"

	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/models"
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
	fmt.Printf("PRINTING USER ROLE: %v", user.Role)

	if user.Role != "teacher" {
		assigh.lg.Println("assignment controller - CreateAssignment - no permission for user")
		return nil
	}
	assigh.lg.Debugln("Creating assignment at controller level")
	tx, err := assigh.DB.Begin(ctx)
	if err != nil {
		assigh.lg.Errorf("assignment controller - CreateAssignment - db begin transaction - %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	var assignmentID int
	// Insert the assignment and get the ID
	err = tx.QueryRow(ctx, "INSERT INTO assignments (name, topicid, weight, teacherid) VALUES ($1, $2, $3, $4) RETURNING id",
		assignment.Name, assignment.TopicID, assignment.Weight, assignment.TeacherID).Scan(&assignmentID)
	if err != nil {
		assigh.lg.Errorf("assignment controller - CreateAssignment - db exec - %v", err)
		return err
	}

	// Link the assignment to students
	for _, studentID := range studentIDs {
		_, err = tx.Exec(ctx, "INSERT INTO assignment_students (assignmentid, studentid) VALUES ($1, $2)",
			assignmentID, studentID)
		if err != nil {
			assigh.lg.Errorf("assignment controller - CreateAssignment - link to students - %v", err)
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		assigh.lg.Errorf("assignment controller - CreateAssignment - commit transaction - %v", err)
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
        SET name = $1, topicid = $2, weight = $3, teacherid = $4
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

func (assignc *AssignmentController) GetAssignmentByID(ctx context.Context, assignmentID int) (models.Assignment, error) {
	var assignment models.Assignment

	err := assignc.DB.QueryRow(ctx, `
        SELECT id, name, topicid, weight, teacherid
        FROM assignments
        WHERE id = $1
    `, assignmentID).Scan(&assignment.ID, &assignment.Name, &assignment.TopicID, &assignment.Weight, &assignment.TeacherID)
	if err != nil {
		assignc.lg.Errorf("assignment controller - GetAssignmentByID - db query - %v", err)
		return models.Assignment{}, err
	}

	return assignment, nil
}

func (assigh *AssignmentController) SubmitAssignment(ctx context.Context, studentID, assignmentID, score int) error {
	assigh.lg.Debugln("Submitting assignment at controller level")
	tx, err := assigh.DB.Begin(ctx)
	if err != nil {
		assigh.lg.Errorf("assignment controller - SubmitAssignment - db begin transaction - %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	// Insert the submission details into the database
	_, err = tx.Exec(ctx, `
        INSERT INTO student_assignments (assignmentid, studentid, score)
        VALUES ($1, $2, $3)
    `, assignmentID, studentID, score)
	if err != nil {
		assigh.lg.Errorf("assignment controller - SubmitAssignment - db exec - %v", err)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		assigh.lg.Errorf("assignment controller - SubmitAssignment - commit transaction - %v", err)
		return err
	}

	return nil
}

func (assignc *AssignmentController) GetCorrectAnswer(ctx context.Context, questionID int) (models.Answer, error) {
	var correctAnswer models.Answer
	err := assignc.DB.QueryRow(ctx, `
        SELECT questionid, text, iscorrect
        FROM answers
        WHERE question_id = $1 AND is_correct = true
    `, questionID).Scan(&correctAnswer.ID, &correctAnswer.Text, &correctAnswer.IsCorrect)
	if err != nil {
		assignc.lg.Errorf("assignment controller - GetCorrectAnswer - db query - %v", err)
		return models.Answer{}, err
	}
	return correctAnswer, nil
}
