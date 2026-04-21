package db

import (
	"testing"
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
