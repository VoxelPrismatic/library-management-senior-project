package db

import (
	"testing"
	"time"
)

func TestLoanFlags(t *testing.T) {
	if val := LoanStatusReturned.ToCopyStatus(); val != CopyLoanAvailable {
		t.Fatalf("loan flags: LoanStatusReturned: expected %d, got %d", CopyLoanAvailable, val)
	}
	if val := LoanStatusCheckedOut.ToCopyStatus(); val != CopyLoanUnvailable {
		t.Fatalf("loan flags: LoanStatusCheckedOut: expected %d, got %d", CopyLoanUnvailable, val)
	}
	if val := LoanStatusOverdue.ToCopyStatus(); val != CopyLoanOverdue {
		t.Fatalf("loan flags: LoanStatusOverdue: expected %d, got %d", CopyLoanOverdue, val)
	}
}

func TestLoanStatusCheckedOut(t *testing.T) {
	loan := Loan{
		DateCheckout: time.Now().Add(-DAY),
		DateReturned: NilTime,
	}

	if loan.Status() != LoanStatusCheckedOut {
		t.Fatalf("loan status: expected LoanStatusCheckedOut, got %s", loan.Status())
	}
}

func TestLoanStatusOverdue(t *testing.T) {
	loan := Loan{
		DateCheckout: time.Now().Add(-LOAN_DURATION + 1),
		DateReturned: NilTime,
	}

	if loan.Status() != LoanStatusOverdue {
		t.Fatalf("loan status: expected LoanStatusOverdue, got %s", loan.Status())
	}
}

func TestLoanStatusReturned(t *testing.T) {
	loan := Loan{
		DateCheckout: time.Now().Add(-WEEK),
		DateReturned: time.Now(),
	}

	if loan.Status() != LoanStatusReturned {
		t.Fatalf("loan status: expected LoanStatusReturned, got %s", loan.Status())
	}
}

func TestLoanReturn(t *testing.T) {
	tx := TestDb()
	defer tx.Rollback()

	user := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
	}
	tx.Save(&user)

	book := BookWork{
		ID:    "test-book-id",
		Title: "Test Book",
	}
	tx.Save(&book)

	c := BookCopy{
		BookWorkID: book.ID,
		Format:     BookFmtPaperback,
		Status:     CopyStatusPublic,
	}
	tx.Save(&c)

	loan := Loan{
		BookCopyID:   c.ID,
		UserID:       user.ID,
		DateCheckout: time.Now().Add(-DAY),
		DateReturned: NilTime,
	}
	tx.Save(&loan)

	err := loan.Return()
	if err != nil {
		t.Fatalf("loan return: unexpected error: %v", err)
	}

	if loan.DateReturned.IsZero() {
		t.Fatalf("loan return: DateReturned should not be zero after return")
	}
}
