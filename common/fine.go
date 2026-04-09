package common

import (
	"time"

	"gorm.io/gorm"
)

type Fine struct {
	gorm.Model
	ID          int `gorm:"primaryKey"`
	UserID      int
	UserRef     User
	LoanID      int
	LoanRef     Loan
	Amount      float32
	PaymentDate time.Time
}
