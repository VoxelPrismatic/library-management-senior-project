package db

import "time"

var _ = Migrate(Hold{})

// When a user wants to get in line to check out a book
type Hold struct {
	BaseModel
	WorkID        SqlUUID `gorm:"type:text"`
	UserID        SqlUUID `gorm:"type:text"`
	RequestDate   time.Time
	FulfilledDate time.Time
	Status        HoldStatus
}

type HoldStatus int

const (
	HoldQueued    HoldStatus = 1 << iota // User in queue
	HoldCancelled                        // User canceled hold
	HoldPostponed                        // User have outstanding charges and cannot check out books right now
	HoldCompleted                        // User has checked out the book
)
