package db

import (
	"time"
)

const (
	DAY           time.Duration = time.Hour * 24
	WEEK                        = DAY * 7
	LOAN_DURATION               = WEEK * 2
)

var _ = Migrate(Loan{})

// Loan record, all checked-out books are like this
// Note: When a continuing hold is satisfied, the previous loan should be
// marked as returned and a new loan should be issued.
type Loan struct {
	BaseModel
	BookCopyID        SqlUUID
	UserID            SqlUUID
	DateCheckout      time.Time
	DateReturned      *time.Time
	OutgoingCondition ConditionFlag
	IncomingCondition ConditionFlag
}

func (l Loan) Overdue() bool {
	if l.DateReturned.IsZero() {
		return false
	}
	return l.DateCheckout.Add(LOAN_DURATION).After(time.Now())
}
