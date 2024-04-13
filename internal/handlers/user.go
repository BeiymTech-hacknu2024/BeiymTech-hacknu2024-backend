package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/controllers"
	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/models"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userc *controllers.UserController
	lg    *logrus.Logger
	store *sessions.CookieStore
}

func NewUserHandler(userc *controllers.UserController, store *sessions.CookieStore, lg *logrus.Logger) *UserHandler {
	return &UserHandler{
		userc: userc,
		store: store,
		lg:    lg,
	}
}

func (userh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	userh.lg.Debugln("User Registration at handler level")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		userh.lg.Error("user handler - register - json decoder - %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = userh.userc.CreateUser(r.Context(), &user)
	if err != nil {
		userh.lg.Error("user handler - register - user create - %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Error(w, "User registered sucessfully", http.StatusCreated)
}

func (userh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	userh.lg.Debugln("User Login at handler level")
	var user models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		userh.lg.Error("user handler - login - json decoder - %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authenticatedUser, err := userh.userc.Authenticate(r.Context(), user.Email, user.Password)
	if err != nil {
		userh.lg.Error("user handler - login - authenticate - %w", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	session, _ := userh.store.Get(r, "session-name")
	session.Values["user_id"] = authenticatedUser.ID
	session.Values["user_email"] = authenticatedUser.Email
	session.Save(r, w)

	http.Error(w, "User logged in sucessfully", http.StatusCreated)
}
