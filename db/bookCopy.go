package db

import (
	"time"
)

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
