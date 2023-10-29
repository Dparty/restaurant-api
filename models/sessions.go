package models

type CreateSessionRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Session struct {
	Token string `json:"token"`
}
