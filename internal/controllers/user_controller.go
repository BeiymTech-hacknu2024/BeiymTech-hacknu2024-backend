package controllers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	m "github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	DB *pgxpool.Pool
	lg *logrus.Logger
}

func (userc *UserController) GetPerformance(ctx context.Context, userID int64) (*m.Performance, error) {
  userc.lg.Debugln("Getting user performance")
  var performance m.Performance

  err := userc.DB.QueryRow(ctx, "SELECT id, kinematics, dynamic, electrodynamics, acids, chemicalbonding, trigonometry, linearalgebra, geometry, probability FROM performance WHERE id = $1", userID).Scan(&performance.ID, &performance.Kinematics, &performance.Dynamics, &performance.Electrodynamics, &performance.Acids, &performance.ChemicalBonding, &performance.Trigonometry, &performance.LinearAlgebra, &performance.Geometry, &performance.Probability)
  if err != nil {
    userc.lg.Errorf("user controller - %v", err)
		return nil, err
  }
  
  return &performance, nil
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

func (userc *UserController) GetUserByID(ctx context.Context, id int) (*m.User, error) {
	userc.lg.Debugln("Getting user by ID at controller level")
	var user m.User
	err := userc.DB.QueryRow(ctx, "SELECT name, email, password, role FROM users WHERE id = $1", id).
		Scan(&user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		userc.lg.Errorf("user controller - GetUserByID - db exec - %v", err)
		return nil, err
	}
	return &user, nil
}

func (userc *UserController) GetUserPerformance(ctx context.Context, userID int) (m.StudentPerformanceBySubject, error) {
	var performance m.StudentPerformanceBySubject

	err := userc.DB.QueryRow(ctx, `
        SELECT subjectid, overallscore, array_agg(assignmentid)
        FROM student_performance_by_subject
        WHERE studentid = $1
        GROUP BY subjectid, overallscore
    `, userID).
		Scan(&performance.SubjectID, &performance.OverallScore, pq.Array(&performance.AssignmentIDList))
	if err != nil {
		userc.lg.Errorf("user controller - GetUserPerformance - db query - %v", err)
		return m.StudentPerformanceBySubject{}, err
	}

	return performance, nil
}
func (userc *UserController) GetAllUsers(ctx context.Context, userID int) ([]m.User, error) {

	users := make([]m.User, 0)

	user, err := userc.GetUserByID(ctx, userID)

	if err != nil {
		userc.lg.Printf("user controller - GetAllUsers - get user by id - %v", err)
		return nil, err
	}
	if user.Role != "teacher" {
		userc.lg.Println("user controller - GetAllUsers - no permission for user")
		return nil, err
	} else {
		rows, err := userc.DB.Query(ctx, "SELECT id, name, email, role FROM users")
		if err != nil {
			userc.lg.Errorf("user controller - GetAllUsers - db query - %v", err)
			return nil, err
		}
		defer rows.Close()
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
	}

	return users, nil
}

func NewUserController(db *pgxpool.Pool, lg *logrus.Logger) *UserController {
	return &UserController{
		DB: db,
		lg: lg,
	}
}
