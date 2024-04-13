package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/controllers"
	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/models"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type AssignmentHandler struct {
	assignc *controllers.AssignmentController
	lg      *logrus.Logger
	store   *sessions.CookieStore
}

func NewAssignmentHandler(assignc *controllers.AssignmentController, store *sessions.CookieStore, lg *logrus.Logger) *AssignmentHandler {
	return &AssignmentHandler{
		assignc: assignc,
		store:   store,
		lg:      lg,
	}
}

func (assignh *AssignmentHandler) CreateAssignment(w http.ResponseWriter, r *http.Request) {

	session, err := assignh.store.Get(r, "session-name")
	if err != nil {
		assignh.lg.Errorf("user handler - GetUser - session get - %v", err)
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Session does not contain user ID", http.StatusUnauthorized)
		return
	}

	assignh.lg.Debugln("Handler level - Create Assignment")
	var req models.CreateAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		assignh.lg.Errorf("assignment handler - CreateAssignment - json decode - %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := assignh.assignc.CreateAssignment(r.Context(), &req.Assignment, req.StudentIDs, userID); err != nil {
		assignh.lg.Errorf("assignment handler - CreateAssignment - %v", err)
		http.Error(w, "Failed to create assignment", http.StatusInternalServerError)
		return
	}

	http.Error(w, "Assignment created successfully", http.StatusCreated)
}

func (assignh *AssignmentHandler) UpdateAssignment(w http.ResponseWriter, r *http.Request) {
	assignh.lg.Debugln("Handler level - Update Assignment")
}
