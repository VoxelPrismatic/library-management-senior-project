package common

import (
	"time"

	"gorm.io/gorm"
)

type Hold struct {
	gorm.Model
	ID          int `gorm:"primaryKey"`
	WorkID      int
	WorkRef     BookWork
	UserID      int
	UserRef     User
	RequestDate time.Time
	Status      HoldStatus
}

type HoldStatus int

const (
	HoldQueued HoldStatus = 1 << iota
	HoldCancelled
	HoldCompleted
)
