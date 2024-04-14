package handlers

import (
	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/controllers"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	UserHandler
	AssignmentHandler
}

func NewHandlers(lg *logrus.Logger, userc *controllers.UserController, assignc *controllers.AssignmentController, store *sessions.CookieStore) *Handlers {
	return &Handlers{
		UserHandler:       *NewUserHandler(userc, store, lg),
		AssignmentHandler: *NewAssignmentHandler(assignc, store, lg),
	}
}
