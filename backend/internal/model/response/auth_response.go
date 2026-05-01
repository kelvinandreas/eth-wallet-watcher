package response

type AuthResponse struct {
	Token string `json:"token"`
}

func NewAuthResponse(token string) AuthResponse {
	return AuthResponse{Token: token}
}
