package db

import "time"

var _ = Migrate(BookCopy{}, RepairLog{})

// An individual copy of a book
type BookCopy struct {
	BaseModel
	BookWork   BookWork
	BookWorkID string
	Barcode    string        // TO-DO: Replace with deterministic function
	Condition  ConditionFlag // TO-DO: Replace with function to derive this based on last return and repair dates
	Format     BookFmtFlag   // Hard-cover, paperback, etc.
	Status     CopyStatusFlag
}

// Repair log for individual copies of a book for audit purposes
type RepairLog struct {
	BaseModel
	BookCopyID     SqlUUID `gorm:"type:text"`
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

type CopyLoanFlag int

const (
	CopyLoanAvailable  CopyLoanFlag = 1 << iota // Open to the public
	CopyLoanUnvailable                          // Book is checked out
	CopyLoanOverdue                             // Book is checked out and overdue
	CopyLoanWithdrawn                           // Book is withdrawn from circulation until repairs are complete
)

type FormatsMap[T any] map[BookFmtFlag]T
type CopyList []BookCopy

func (arr CopyList) MapFormats() FormatsMap[CopyList] {
	ret := FormatsMap[CopyList]{}
	for _, e := range arr {
		_, exists := ret[e.Format]
		if !exists {
			ret[e.Format] = CopyList{e}
		} else {
			ret[e.Format] = append(ret[e.Format], e)
		}

	}
	return ret
}

func (c BookCopy) LoanHistory() ([]Loan, error) {
	db := Db()
	ret := []Loan{}
	status := db.Model(&Loan{}).
		Where(&Loan{
			BookCopyID: c.ID,
		}).
		Order("date_checkout DESC").
		Preload("User").
		Preload("BookCopy").
		Find(&ret)
	return ret, status.Error
}

func (c BookCopy) LoanStatus() (CopyLoanFlag, error) {
	if c.Status != CopyStatusPublic {
		return CopyLoanWithdrawn, nil
	}

	history, err := c.LoanHistory()
	if err != nil {
		return CopyLoanWithdrawn, err
	}

	if len(history) == 0 {
		return CopyLoanAvailable, nil
	}

	return history[0].Status().ToCopyStatus(), nil
}
