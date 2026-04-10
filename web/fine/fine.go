package fine

import (
	"time"

	"voxelprismatic/library-management-senior-project/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ = db.Migrate(Fine{})

// A fine that must be paid before the user can check out more books.
// A fine can be issued for any reason, including late, lost, or damaged.
type Fine struct {
	gorm.Model
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID
	LoanID uuid.UUID

	IssueReason     FineReasonFlag
	IssueDate       time.Time
	AmountIssued    float32
	AmountRemaining float32 // How much is remaining, recalculated after every transaction

	AmountWaived  float32 // Any discounts provided by a librarian
	WaivedReasion string
	WaivedBy      uuid.UUID
}

type FineReasonFlag int

const (
	FineReasonLate    FineReasonFlag = iota // Did not return the book on time
	FineReasonLost                          // Lost the book; fee for replacement
	FineReasonDamaged                       // Book was received damaged, eg torn pages
)
