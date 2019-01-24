package main

import (
	"github.com/jinzhu/gorm"
)

// User model
type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(255);unique;not null"`
}

// Group model
type Group struct {
	gorm.Model
	Name        string        `gorm:"type:varchar(255);unique;not null"`
	Permissions []*Permission `gorm:"many2many:group_permissions;"`
	Users       []*User       `gorm:"many2many:group_users;"`
}

// Permission model
type Permission struct {
	gorm.Model
	Name string `gorm:"type:varchar(255);unique;not null"`
}
