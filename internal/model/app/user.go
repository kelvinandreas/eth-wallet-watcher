package app

type User struct {
	BaseModel
	Email    string `gorm:"size:255;not null;uniqueIndex"`
	Password string `gorm:"not null"`
}
