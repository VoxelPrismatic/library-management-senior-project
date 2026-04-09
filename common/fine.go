package common

import (
	"time"

	"gorm.io/gorm"
)

type Fine struct {
	gorm.Model
	ID             int `gorm:"primaryKey"`
	UserID         int
	UserRef        User
	LoanID         int
	LoanRef        Loan
	AmountAssessed float32
	AmountPaid     float32
	Status         FineStatus
	LastPayment    time.Time
}

type FineStatus int

const (
	FineStatusUnpaid FineStatus = iota
	FineStatusPartial
	FineStatusPaid
	FineStatusWaived
)
