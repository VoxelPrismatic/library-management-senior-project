package common

import (
	"time"

	"gorm.io/gorm"
)

type BookCopy struct {
	gorm.Model
	ID        int `gorm:"primaryKey"`
	WorkID    int
	WorkRef   BookWork
	Barcode   string
	Condition ConditionFlag
	Status    CopyStatusFlag
}

type ConditionFlag int

const (
	ConditionMint ConditionFlag = iota // New from the factory
	ConditionGood                      // No major wear, but some pages are bent
	ConditionFair                      // Light wear on corners, crease marks, but no torn pages
	ConditionPoor                      // Some tears, annotations, etc
	ConditionDead                      // Missing pages
)

type CopyStatusFlag int

const (
	CopyStatusPublic        CopyStatusFlag = iota // Open to the public
	CopyStatusPendingReturn                       // Book is checked out, waiting to be removed
	CopyStatusPendingAction                       // Book is returned, but waiting for action (repair or discard)
	CopyStatusRepairing                           // Book is being repaired (rebound, etc.)
	CopyStatusDiscarded                           // Book is discarded, possibly replaced
)

type BookCopyRepairLog struct {
	gorm.Model
	ID             int `gorm:"primaryKey"`
	BookCopyID     int
	BookCopyRef    BookCopy
	Date           time.Time
	IncomingStatus CopyStatusFlag
	OutgoingStatus CopyStatusFlag
	Reason         string
	AuthorizedBy   string
}
