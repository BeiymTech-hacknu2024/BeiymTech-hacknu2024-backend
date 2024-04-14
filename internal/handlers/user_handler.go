package handlers

import (
	"encoding/json"
	"fmt"
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

func (userh *UserHandler) GeneratePerformance(w http.ResponseWriter, r *http.Request) {
	userh.lg.Debugln("User Performance at handler level")
	var performance *models.Performance // Change type to pointer
	session, err := userh.store.Get(r, "session-name")
	fmt.Println(session)
	if err != nil {
		userh.lg.Errorf("user handler - GeneratePerformance - session get - %v", err)
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Session does not contain user ID", http.StatusUnauthorized)
		return
	}
	performance, err = userh.userc.GetPerformance(r.Context(), userID) // Remove := to use the already declared performance variable
	if err != nil {
		http.Error(w, "Failed to retrieve user performance", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(performance)
	if err != nil {
		userh.lg.Errorf("user handler - GetUser - json marshal - %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
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

func (userh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	session, err := userh.store.Get(r, "session-name")
	if err != nil {
		userh.lg.Errorf("user handler - GetUser - session get - %v", err)
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Session does not contain user ID", http.StatusUnauthorized)
		return
	}

	user, err := userh.userc.GetUserByID(r.Context(), userID)
	if err != nil {
		userh.lg.Errorf("user handler - GetUser - get user by id - %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(user)
	if err != nil {
		userh.lg.Errorf("user handler - GetUser - json marshal - %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (userh *UserHandler) GetUserPerformance(w http.ResponseWriter, r *http.Request) {
	// Get user ID from session
	session, err := userh.store.Get(r, "session-name")
	if err != nil {
		userh.lg.Errorf("user handler - GetUserPerformance - session get - %v", err)
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Session does not contain user ID", http.StatusUnauthorized)
		return
	}

	// Fetch user's performance data from the controller
	performance, err := userh.userc.GetUserPerformance(r.Context(), userID)
	if err != nil {
		userh.lg.Errorf("user handler - GetUserPerformance - get user performance - %v", err)
		http.Error(w, "Failed to retrieve user performance", http.StatusInternalServerError)
		return
	}

	// Convert performance data to JSON
	jsonResponse, err := json.Marshal(performance)
	if err != nil {
		userh.lg.Errorf("user handler - GetUserPerformance - json marshal - %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (userh *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	session, err := userh.store.Get(r, "session-name")
	if err != nil {
		userh.lg.Errorf("user handler - GetUser - session get - %v", err)
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Session does not contain user ID", http.StatusUnauthorized)
		return
	}

	users, err := userh.userc.GetAllUsers(r.Context(), userID)
	if err != nil {
		userh.lg.Errorf("user handler - GetUsers - get all users - %v", err)
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(users)
	if err != nil {
		userh.lg.Errorf("user handler - GetUsers - json marshal - %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (userh *UserHandler) GetAllPerformance(w http.ResponseWriter, r *http.Request) {
	session, err := userh.store.Get(r, "session-name")
	if err != nil {
		userh.lg.Errorf("user handler - GetUser - session get - %v", err)
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Session does not contain user ID", http.StatusUnauthorized)
		return
	}

	performances, err := userh.userc.GetAllUserPerformance(r.Context(), userID)
	if err != nil {
		userh.lg.Errorf("user handler - GetUsers - get all performances - %v", err)
		http.Error(w, "Failed to retrieve performances", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(performances)
	if err != nil {
		userh.lg.Errorf("user handler - GetUsers - json marshal - %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
