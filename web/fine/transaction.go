package fine

import (
	"time"
	"voxelprismatic/library-management-senior-project/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ = db.Migrate(Transaction{})

// Transactions resolve the oldest fines first
/* TO-DO: Figure out how to link transactions to specific fines
 *        - Do we split one transaction into its components?
 *        - Do we not care?
 *        - Do we include a list of Fine IDs?
 */
type Transaction struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID
	AmountPaid float32
	Date       time.Time
}
