package book

import (
	"time"
	"voxelprismatic/library-management-senior-project/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ = db.Migrate(BookCopy{}, RepairLog{})

// An individual copy of a book
type BookCopy struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	WorkID    string
	Barcode   string        // TO-DO: Replace with deterministic function
	Condition ConditionFlag // TO-DO: Replace with function to derive this based on last return and repair dates
	Format    BookFmtFlag   // Hard-cover, paperback, etc.
	Status    CopyStatusFlag
}

// Repair log for individual copies of a book for audit purposes
type RepairLog struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	BookCopyID     uuid.UUID
	Date           time.Time
	IncomingStatus CopyStatusFlag
	OutgoingStatus CopyStatusFlag
	TechnicianName string
}

type BookFmtFlag int

const (
	BookFmtPaperback     BookFmtFlag = 1 << iota // Physical paperback book
	BookFmtHardCover                             // Physical hard-cover book
	BookFmtPhysicalAudio                         // E.g. a physical MP3 player with the book preloaded
	BookFmtDigitalBook                           // E.g. Kindle
	BookFmtDigitalAudio                          // E.g. Audible
)

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
