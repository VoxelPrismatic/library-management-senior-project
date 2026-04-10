package fetch

import (
	"fmt"
	"testing"
	"voxelprismatic/library-management-senior-project/web/book"
)

func TestLookup(t *testing.T) {
	isbn := 9781613163528 // Andrew Klavan, A Strange Habit of Mind
	data, err := GBooksIsbnLookup(isbn)
	if err != nil {
		t.Errorf(`GBooksIsbnLookup(isbn) unexpectedly returned error %v`, err)
		return
	}

	if len(data.Items) == 0 {
		t.Errorf(`GBooksIsbnLookup(isbn) yield zero results, expected at least one`)
		return
	}

	found := false
	isbn_str := fmt.Sprint(isbn)
	for _, obj := range data.Items {
		for _, id := range obj.VolumeInfo.IndustryIdentifiers {
			if id.Identifier == isbn_str {
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		t.Errorf(`GBooksIsbnLookup(isbn) did not find isbn %d`, isbn)
		return
	}
}

func TestVolume(t *testing.T) {
	vol := book.GBooksVolumeID_t("Bj6VEAAAQBAJ")
	data, err := GBooksVolume(vol)
	if err != nil {
		t.Errorf(`GBooksVolume(vol) unexpectedly returned error %v`, err)
		return
	}

	if data.ID != string(vol) {
		t.Errorf(`GBooksVolume(vol) expected "%s", got "%s"`, vol, data.ID)
		return
	}
}
