package models

type User struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Password string `json:"password"`
}



type LoginRequest struct {
  Email    string `json:"email"`
	Password string `json:"password"`
}

