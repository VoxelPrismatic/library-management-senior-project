package db

import "time"

var _ = Migrate(Hold{})

// When a user wants to get in line to check out a book
type Hold struct {
	BaseModel
	BookWork      BookWork
	BookWorkID    SqlUUID `gorm:"type:text"`
	User          User
	UserID        SqlUUID `gorm:"type:text"`
	RequestedDate time.Time
	FulfilledDate time.Time
	CancelledDate time.Time
}

type HoldStatus int

const (
	HoldQueued    HoldStatus = 1 << iota // User in queue
	HoldCancelled                        // User canceled hold
	HoldPostponed                        // User have outstanding charges and cannot check out books right now
	HoldCompleted                        // User has checked out the book
	HoldRevoked                          // User account has been deleted
)

func (h Hold) Status() (HoldStatus, error) {
	if !h.FulfilledDate.IsZero() {
		return HoldCompleted, nil
	}

	if !h.CancelledDate.IsZero() {
		return HoldCancelled, nil
	}

	u := h.GetUser()
	if u.Status == UserStatusDeleted {
		return HoldRevoked, nil
	}

	if state, err := u.HasOverdueBooks(); err != nil {
		return 0, err
	} else if state {
		return HoldPostponed, nil
	}

	if count, err := u.CheckedOut(); err != nil {
		return 0, err
	} else if len(count) >= LOAN_LIMIT {
		return HoldPostponed, nil
	}

	return HoldQueued, nil
}

func (h Hold) GetUser() User {
	if h.User.ID.IsEmpty() {
		db.Where(&User{BaseModel: BaseModel{ID: h.UserID}}).First(&h.User)
	}
	return h.User
}
