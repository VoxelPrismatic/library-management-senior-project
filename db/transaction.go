package db

import "time"

var _ = Migrate(Transaction{})

// Transactions resolve the oldest fines first
/* TO-DO: Figure out how to link transactions to specific fines
 *        - Do we split one transaction into its components?
 *        - Do we not care?
 *        - Do we include a list of Fine IDs?
 */
type Transaction struct {
	BaseModel
	UserID     SqlUUID `gorm:"type:text"`
	AmountPaid float32
	Date       time.Time
}
