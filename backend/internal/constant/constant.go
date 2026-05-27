package constant

import "github.com/google/uuid"

var SystemUserID = uuid.Nil

const (
	// JWT errors
	ErrMissingAuthorizationHeader = "Missing authorization header"
	ErrInvalidAuthorizationFormat = "Invalid authorization format"
	ErrInvalidOrExpiredToken      = "Invalid or expired token"
	ErrInvalidTokenClaims         = "Invalid token claims"

	// Auth errors
	ErrEmailAlreadyExists       = "Email already exists"
	ErrFailedHashPassword       = "Failed to hash password"
	ErrInvalidEmailOrPassword   = "Invalid email or password"
	ErrFailedGenerateToken      = "Failed to generate token"
	ErrInvalidRequestBody       = "Invalid request body"
	ErrEmailAndPasswordRequired = "Email and password are required"

	// Auth messages
	MsgRegistrationSuccessful = "Registration successful"
	MsgLoginSuccessful        = "Login successful"

	// Wallet errors
	ErrWalletAddressRequired = "Wallet address is required"
	ErrInvalidEthAddress     = "Invalid Ethereum address format"
	ErrWalletAlreadyExists   = "Wallet address already added"
	ErrInvalidUserID         = "Invalid user ID"
	ErrInvalidWalletID       = "Invalid wallet ID"
	ErrWalletNotFound        = "Wallet not found"

	// Wallet messages
	MsgWalletCreated   = "Wallet created successfully"
	MsgWalletRetrieved = "Wallets retrieved successfully"
	MsgWalletDeleted   = "Wallet deleted successfully"

	// Transaction messages
	MsgTransactionsRetrieved = "Transactions retrieved successfully"

	// Notification errors
	ErrInvalidNotificationID = "Invalid notification ID"
	ErrNotificationNotFound  = "Notification not found"

	// Notification messages
	MsgNotificationsRetrieved = "Notifications retrieved successfully"
	MsgNotificationRead       = "Notification marked as read"
)

var CacheKey = struct {
	Transactions  string
	Notifications string
}{
	Transactions:  "transactions:%s",
	Notifications: "notifications:%s",
}
