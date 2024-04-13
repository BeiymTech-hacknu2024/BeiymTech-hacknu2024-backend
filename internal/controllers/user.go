package controllers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	m "github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	DB *pgxpool.Pool
	lg *logrus.Logger
}

func NewUserController(db *pgxpool.Pool, lg *logrus.Logger) *UserController {
	return &UserController{
		DB: db,
		lg: lg,
	}
}

func (userc *UserController) CreateUser(ctx context.Context, user *m.User) error {
	userc.lg.Debugln("User Creation at controller level")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		userc.lg.Errorf("user controller - CreateUser - password generation - %v", err)
		return err
	}
	_, err = userc.DB.Exec(ctx,
		"INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4)",
		user.Name, user.Email, string(hashedPassword), user.Role)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") && strings.Contains(err.Error(), "Email") {
			return fmt.Errorf("email %s already exists", user.Email)
		}
		userc.lg.Errorf("user controller - CreateUser - db exec - %v", err)
		return err
	}
	return nil
}

func (userc *UserController) Authenticate(ctx context.Context, email, password string) (m.User, error) {
	userc.lg.Debugln("User Authentication at controller level")
	var user m.User
	err := userc.DB.QueryRow(ctx, "SELECT id, email, password FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		userc.lg.Errorf("user controller - Authenticate - db exec - %v", err)
		return m.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		userc.lg.Errorf("user controller - Authenticate - hash and password comparison - %v", err)
		return m.User{}, errors.New("incorrect password")
	}
	return user, nil
}

func (userc *UserController) GetUserByID(ctx context.Context, id int64) (*m.User, error) {
	userc.lg.Debugln("Getting user by ID at controller level")
	var user m.User
	err := userc.DB.QueryRow(ctx, "SELECT id, name, role, passwrod FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		userc.lg.Errorf("user controller - GetUserByID - db exec - %v", err)
		return nil, err
	}
	return &user, nil
}

func (userc *UserController) CreateAssignment(ctx context.Context, assignment *m.Assignment, studentIDs []int) error {
	userc.lg.Debugln("Creating assignment at controller level")
	tx, err := userc.DB.Begin(ctx)
	if err != nil {
		userc.lg.Errorf("user controller - CreateAssignment - db begin transaction - %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	var assignmentID int
	// Insert the assignment and get the ID
	err = tx.QueryRow(ctx, "INSERT INTO assignments (name, topic_id, weight, teacher_id) VALUES ($1, $2, $3, $4) RETURNING id",
		assignment.Name, assignment.TopicID, assignment.Weight, assignment.TeacherID).Scan(&assignmentID)
	if err != nil {
		userc.lg.Errorf("user controller - CreateAssignment - db exec - %v", err)
		return err
	}

	// Link the assignment to students
	for _, studentID := range studentIDs {
		_, err = tx.Exec(ctx, "INSERT INTO assignment_students (assignment_id, student_id) VALUES ($1, $2)",
			assignmentID, studentID)
		if err != nil {
			userc.lg.Errorf("user controller - CreateAssignment - link to students - %v", err)
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		userc.lg.Errorf("user controller - CreateAssignment - commit transaction - %v", err)
		return err
	}

	return nil
}

func (userc *UserController) GetAllUsers(ctx context.Context) ([]m.User, error) {
	userc.lg.Debugln("Getting all users at controller level")

	rows, err := userc.DB.Query(ctx, "SELECT id, name, email, role FROM users")
	if err != nil {
		userc.lg.Errorf("user controller - GetAllUsers - db query - %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []m.User
	for rows.Next() {
		var user m.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role); err != nil {
			userc.lg.Errorf("user controller - GetAllUsers - scan row - %v", err)
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		userc.lg.Errorf("user controller - GetAllUsers - rows iteration - %v", err)
		return nil, err
	}

	return users, nil
}
