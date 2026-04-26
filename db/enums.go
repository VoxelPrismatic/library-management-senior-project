package db

import "fmt"

type BookFmtFlag int

const (
	BookFmtPaperback     BookFmtFlag = 1 << iota // Physical paperback book
	BookFmtHardCover                             // Physical hard-cover book
	BookFmtPhysicalAudio                         // E.g. a physical MP3 player with the book preloaded
	BookFmtDigitalBook                           // E.g. Kindle
	BookFmtDigitalAudio                          // E.g. Audible
)

func (f BookFmtFlag) String() string {
	switch f {
	case BookFmtPaperback:
		return "BookFmtPaperback"
	case BookFmtHardCover:
		return "BookFmtHardCover"
	case BookFmtPhysicalAudio:
		return "BookFmtPhysicalAudio"
	case BookFmtDigitalBook:
		return "BookFmtDigitalBook"
	case BookFmtDigitalAudio:
		return "BookFmtDigitalAudio"
	default:
		return fmt.Sprintf("BookFmtFlag(%d)", f)
	}
}

type ConditionFlag int

const (
	ConditionMint ConditionFlag = iota // New from the factory
	ConditionGood                      // No major wear, but some pages are bent
	ConditionFair                      // Light wear on corners, crease marks, but no torn pages
	ConditionPoor                      // Some tears, annotations, etc
	ConditionDead                      // Missing pages
)

func (f ConditionFlag) String() string {
	switch f {
	case ConditionMint:
		return "ConditionMint"
	case ConditionGood:
		return "ConditionGood"
	case ConditionFair:
		return "ConditionFair"
	case ConditionPoor:
		return "ConditionPoor"
	case ConditionDead:
		return "ConditionDead"
	default:
		return fmt.Sprintf("ConditionFlag(%d)", f)
	}
}

type CopyStatusFlag int

const (
	CopyStatusPublic        CopyStatusFlag = iota // Open to the public
	CopyStatusPendingReturn                       // Book is checked out, waiting to be removed
	CopyStatusPendingAction                       // Book is returned, but waiting for action (repair or discard)
	CopyStatusRepairing                           // Book is being repaired (rebound, etc.)
	CopyStatusDiscarded                           // Book is discarded, possibly replaced
)

func (f CopyStatusFlag) String() string {
	switch f {
	case CopyStatusPublic:
		return "CopyStatusPublic"
	case CopyStatusPendingReturn:
		return "CopyStatusPendingReturn"
	case CopyStatusPendingAction:
		return "CopyStatusPendingAction"
	case CopyStatusRepairing:
		return "CopyStatusRepairing"
	case CopyStatusDiscarded:
		return "CopyStatusDiscarded"
	default:
		return fmt.Sprintf("CopyStatusFlag(%d)", f)
	}
}

type CopyLoanFlag int

const (
	CopyLoanAvailable  CopyLoanFlag = 1 << iota // Open to the public
	CopyLoanUnvailable                          // Book is checked out
	CopyLoanOverdue                             // Book is checked out and overdue
	CopyLoanWithdrawn                           // Book is withdrawn from circulation until repairs are complete
)

func (f CopyLoanFlag) String() string {
	switch f {
	case CopyLoanAvailable:
		return "CopyLoanAvailable"
	case CopyLoanUnvailable:
		return "CopyLoanUnvailable"
	case CopyLoanOverdue:
		return "CopyLoanOverdue"
	case CopyLoanWithdrawn:
		return "CopyLoanWithdrawn"
	default:
		return fmt.Sprintf("CopyLoanFlag(%d)", f)
	}
}

type FineReasonFlag int

const (
	FineReasonLate    FineReasonFlag = iota // Did not return the book on time
	FineReasonLost                          // Lost the book; fee for replacement
	FineReasonDamaged                       // Book was received damaged, eg torn pages
)

func (f FineReasonFlag) String() string {
	switch f {
	case FineReasonLate:
		return "FineReasonLate"
	case FineReasonLost:
		return "FineReasonLost"
	case FineReasonDamaged:
		return "FineReasonDamaged"
	default:
		return fmt.Sprintf("FineReasonFlag(%d)", f)
	}
}

type HoldStatusFlag int

const (
	HoldQueued    HoldStatusFlag = 1 << iota // User in queue
	HoldCancelled                            // User canceled hold
	HoldPostponed                            // User have outstanding charges and cannot check out books right now
	HoldCompleted                            // User has checked out the book
	HoldRevoked                              // User account has been deleted
)

func (f HoldStatusFlag) String() string {
	switch f {
	case HoldQueued:
		return "HoldQueued"
	case HoldCancelled:
		return "HoldCancelled"
	case HoldPostponed:
		return "HoldPostponed"
	case HoldCompleted:
		return "HoldCompleted"
	case HoldRevoked:
		return "HoldRevoked"
	default:
		return fmt.Sprintf("HoldStatus(%d)", f)
	}
}

type LoanStatusFlag int

const (
	LoanStatusReturned   LoanStatusFlag = 1 << iota // Book has been returned
	LoanStatusCheckedOut                            // Book is currently checked out
	LoanStatusOverdue                               // Book is checked out, but overdue
)

func (s LoanStatusFlag) String() string {
	switch s {
	case LoanStatusReturned:
		return "LoanStatusReturned"
	case LoanStatusCheckedOut:
		return "LoanStatusCheckedOut"
	case LoanStatusOverdue:
		return "LoanStatusOverdue"
	default:
		return fmt.Sprintf("LoanStatusFlag(%d)", s)
	}
}

type UserRoleFlag int

const (
	UserRoleNone      UserRoleFlag = 0         // Logged out user
	UserRolePublic    UserRoleFlag = 1 << iota // General public user
	UserRoleLibrarian                          // Librarian
	UserRoleAdmin                              // Administrator
)

func (f UserRoleFlag) String() string {
	switch f {
	case UserRoleNone:
		return "UserRoleNone"
	case UserRolePublic:
		return "UserRolePublic"
	case UserRoleLibrarian:
		return "UserRoleLibrarian"
	case UserRoleAdmin:
		return "UserRoleAdmin"
	default:
		return fmt.Sprintf("UserRoleFlag(%d)", f)
	}
}

type UserStatusFlag int

const (
	UserStatusActive  UserStatusFlag = 1 << iota // No outstanding issues
	UserStatusLimited                            // User has hit loan limit
	UserStatusLocked                             // TO-DO: Check if user has outstanding fees and remove this redundant lock
	UserStatusDeleted                            // For audit purposes; we may choose to anonymize any data
)

func (f UserStatusFlag) String() string {
	switch f {
	case UserStatusActive:
		return "UserStatusActive"
	case UserStatusLimited:
		return "UserStatusLimited"
	case UserStatusLocked:
		return "UserStatusLocked"
	case UserStatusDeleted:
		return "UserStatusDeleted"
	default:
		return fmt.Sprintf("UserStatusFlag(%d)", f)
	}
}
