package routes

import (
	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(userh *handlers.UserHandler, assignh *handlers.AssignmentHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/register", userh.Register)
	r.Post("/login", userh.Login)

	r.With(userh.RequireAuth).Get("/users", userh.GetUsers)
	r.With(userh.RequireAuth).Get("/user", userh.GetUser)

	r.With(userh.RequireAuth).Post("/assignments/create", assignh.CreateAssignment)
	r.With(userh.RequireAuth).Patch("/assignment/update", assignh.UpdateAssignment)

	return r

}
