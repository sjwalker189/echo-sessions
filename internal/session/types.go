package session

type AuthUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type FlashMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
