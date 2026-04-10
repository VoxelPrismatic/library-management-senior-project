package loan

import (
	"time"

	"voxelprismatic/library-management-senior-project/db"
	"voxelprismatic/library-management-senior-project/web/book"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	DAY           time.Duration = time.Hour * 24
	WEEK                        = DAY * 7
	LOAN_DURATION               = WEEK * 2
)

var _ = db.Migrate(Loan{})

// Loan record, all checked-out books are like this
// Note: When a continuing hold is satisfied, the previous loan should be
// marked as returned and a new loan should be issued.
type Loan struct {
	gorm.Model
	ID                uuid.UUID `gorm:"type:uuid;primaryKey"`
	BookCopyID        uuid.UUID
	UserID            uuid.UUID
	DateCheckout      time.Time
	DateReturned      *time.Time
	OutgoingCondition book.ConditionFlag
	IncomingCondition book.ConditionFlag
}

func (l Loan) Overdue() bool {
	if l.DateReturned.Equal(db.NilTime) {
		return false
	}
	return l.DateCheckout.Add(LOAN_DURATION).After(time.Now())
}
