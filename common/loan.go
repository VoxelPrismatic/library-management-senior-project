package common

import (
	"time"

	"gorm.io/gorm"
)

type Loan struct {
	gorm.Model
	ID           int `gorm:"primaryKey"`
	BookCopyID   int
	BookCopyRef  BookCopy
	UserID       int
	UserRef      User
	DateCheckout time.Time
	DateDue      time.Time
	DateReturned time.Time
}
