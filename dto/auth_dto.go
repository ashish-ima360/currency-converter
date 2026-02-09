package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	UserId  int    `json:"user_id"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

type LoginResult struct {
	ID    int
	Token string
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`	
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterResponse struct {
	UserId  int    `json:"user_id"`
	Message string `json:"message"`
}

