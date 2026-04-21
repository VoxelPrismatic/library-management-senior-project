package db

import (
	"time"
)

const (
	DAY           time.Duration = time.Hour * 24
	WEEK                        = DAY * 7
	LOAN_DURATION               = WEEK * 2 // Maximum amount of time before a book is considered overdue
	LOAN_LIMIT                  = 8        // Maximum amount of books that can be checked out at once
)

var _ = Migrate(Loan{})

// Loan record, all checked-out books are like this
// Note: When a continuing hold is satisfied, the previous loan should be
// marked as returned and a new loan should be issued.
type Loan struct {
	BaseModel
	BookCopy          BookCopy
	BookCopyID        SqlUUID `gorm:"type:text"`
	User              User
	UserID            SqlUUID `gorm:"type:text"`
	DateCheckout      time.Time
	DateReturned      time.Time
	OutgoingCondition ConditionFlag
	IncomingCondition ConditionFlag
}

type LoanStatusFlag int

const (
	LoanStatusReturned LoanStatusFlag = 1 << iota
	LoanStatusCheckedOut
	LoanStatusOverdue
)

func (s LoanStatusFlag) ToCopyStatus() CopyLoanFlag {
	switch s {
	case LoanStatusReturned:
		return CopyLoanAvailable
	case LoanStatusCheckedOut:
		return CopyLoanUnvailable
	case LoanStatusOverdue:
		return CopyLoanOverdue
	default:
		panic("unreachable")
	}
}

func (l Loan) Status() LoanStatusFlag {
	if !l.DateReturned.IsZero() {
		return LoanStatusReturned
	}
	if l.DateReturned.Add(LOAN_DURATION).Before(time.Now()) {
		return LoanStatusOverdue
	}
	return LoanStatusCheckedOut
}

// Marks a book as returned
func (c *Loan) Return() error {
	c.DateReturned = time.Now()
	return db.Save(c).Error
}
