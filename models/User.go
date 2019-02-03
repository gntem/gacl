package models

// User model
type User struct {
	ID   int64
	Name string `gorm:"type:varchar(255);unique;not null"`
}
