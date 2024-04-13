package routes

import (
	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(userh *handlers.UserHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/register", userh.Register)
	r.Post("/login", userh.Login)

	return r
	


}


