package db

import (
	"testing"
	"time"
)

func TestLoanStatusCheckedOut(t *testing.T) {
	setupTestDb(t)

	loan := Loan{
		DateCheckout: time.Now().Add(-DAY),
		DateReturned: NilTime,
	}

	if loan.Status() != LoanStatusCheckedOut {
		t.Fatalf("loan status: expected LoanStatusCheckedOut, got %d", loan.Status())
	}
}

func TestLoanStatusOverdue(t *testing.T) {
	setupTestDb(t)

	loan := Loan{
		DateCheckout: time.Now().Add(-LOAN_DURATION * 2),
		DateReturned: NilTime,
	}

	if loan.Status() != LoanStatusOverdue {
		t.Fatalf("loan status: expected LoanStatusOverdue, got %d", loan.Status())
	}
}

func TestLoanStatusReturned(t *testing.T) {
	setupTestDb(t)

	loan := Loan{
		DateCheckout: time.Now().Add(-WEEK),
		DateReturned: time.Now(),
	}

	if loan.Status() != LoanStatusReturned {
		t.Fatalf("loan status: expected LoanStatusReturned, got %d", loan.Status())
	}
}

func TestLoanReturn(t *testing.T) {
	setupTestDb(t)

	user := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
	}
	db.Save(&user)

	book := BookWork{
		ID:    "test-book-id",
		Title: "Test Book",
	}
	db.Save(&book)

	copy := BookCopy{
		BookWorkID: book.ID,
		Format:     BookFmtPaperback,
		Status:     CopyStatusPublic,
	}
	db.Save(&copy)

	loan := Loan{
		BookCopyID:   copy.ID,
		UserID:       user.ID,
		DateCheckout: time.Now().Add(-DAY),
		DateReturned: NilTime,
	}
	db.Save(&loan)

	err := loan.Return()
	if err != nil {
		t.Fatalf("loan return: unexpected error: %v", err)
	}

	if loan.DateReturned.IsZero() {
		t.Fatalf("loan return: DateReturned should not be zero after return")
	}
}
