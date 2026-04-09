package common

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        int `gorm:"primaryKey"`
	Roles     UserRoleFlag
	FirstName string
	LastName  string
	Email     string
	Status    UserStatusFlag
}

type UserRoleFlag int

const (
	UserRolePublic UserRoleFlag = 1 << iota
	UserRoleAdmin
)

type UserStatusFlag int

const (
	UserStatusActive UserStatusFlag = 1 << iota
	UserStatusLocked                // E.g. user has too many fines and must pay before checking out new books
	UserStatusDeleted
)
