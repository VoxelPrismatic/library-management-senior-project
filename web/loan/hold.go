package loan

import (
	"time"

	"voxelprismatic/library-management-senior-project/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ = db.Migrate(Hold{})

// When a user wants to get in line to check out a book
type Hold struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	WorkID      string
	UserID      uuid.UUID
	RequestDate time.Time
	Status      HoldStatus
}

type HoldStatus int

const (
	HoldQueued    HoldStatus = 1 << iota // User in queue
	HoldCancelled                        // User canceled hold
	HoldPostponed                        // User have outstanding charges and cannot check out books right now
	HoldCompleted                        // User has checked out the book
)
