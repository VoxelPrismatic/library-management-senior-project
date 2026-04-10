package common

import (
	"time"

	"gorm.io/gorm"
)

// A fine that must be paid before the user can check out more books.
// A fine can be issued for any reason, including late, lost, or damaged.
type Fine struct {
	gorm.Model
	ID      int `gorm:"primaryKey"`
	UserID  int
	UserRef User
	LoanID  int
	LoanRef Loan

	IssueReason     FineReasonFlag
	IssueDate       time.Time
	AmountIssued    float32
	AmountRemaining float32 // How much is remaining, recalculated after every transaction

	AmountWaived  float32 // Any discounts provided by a librarian
	WaivedReasion string
	WaivedBy      string
}

// Transactions resolve the oldest fines first
/* TO-DO: Figure out how to link transactions to specific fines
 *        - Do we split one transaction into its components?
 *        - Do we not care?
 *        - Do we include a list of Fine IDs?
 */
type Transaction struct {
	gorm.Model
	UserID     int
	UserRef    User
	AmountPaid float32
	Date       time.Time
}

type FineReasonFlag int

const (
	FineReasonLate FineReasonFlag = iota
	FineReasonLost
	FineReasonDamaged
)
