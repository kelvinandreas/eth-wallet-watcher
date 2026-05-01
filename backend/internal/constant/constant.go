package constant

import "github.com/google/uuid"

var SystemUserID = uuid.Nil

const (
	// Auth errors
	ErrEmailAlreadyExists       = "Email already exists"
	ErrFailedHashPassword       = "Failed to hash password"
	ErrInvalidEmailOrPassword   = "Invalid email or password"
	ErrFailedGenerateToken      = "Failed to generate token"
	ErrInvalidRequestBody       = "Invalid request body"
	ErrEmailAndPasswordRequired = "Email and password are required"

	// JWT errors
	ErrMissingAuthorizationHeader = "Missing authorization header"
	ErrInvalidAuthorizationFormat = "Invalid authorization format"
	ErrInvalidOrExpiredToken      = "Invalid or expired token"
	ErrInvalidTokenClaims         = "Invalid token claims"

	// Auth messages
	MsgRegistrationSuccessful = "Registration successful"
	MsgLoginSuccessful        = "Login successful"
)
