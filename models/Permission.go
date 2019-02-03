package models

// Permission model
type Permission struct {
	ID   int64
	Name string `gorm:"type:varchar(255);unique;not null"`
}
