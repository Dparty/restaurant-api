package models

type CreateSessionRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Session struct {
	Token string `json:"token"`
}

type Account struct {
}

type PutRestaurantRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
