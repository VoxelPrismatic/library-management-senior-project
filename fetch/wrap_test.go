package fetch

import (
	"fmt"
	"testing"
)

func TestLookup(t *testing.T) {
	isbn := "9781613163528" // Andrew Klavan, A Strange Habit of Mind
	data, err := GBooksSearch(fmt.Sprintf("isbn:%s", isbn))
	if err != nil {
		t.Errorf(`unexpectedly returned error %v`, err)
		return
	}

	if len(data.Items) == 0 {
		t.Errorf(`yield zero results, expected at least one`)
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
		t.Errorf(`did not find isbn %s`, isbn)
		return
	}
}

func TestVolume(t *testing.T) {
	vol := "Bj6VEAAAQBAJ"
	data, err := GBooksVolume(vol)
	if err != nil {
		t.Errorf(`unexpectedly returned error %v`, err)
		return
	}

	if data.ID != string(vol) {
		t.Errorf(`expected "%s", got "%s"`, vol, data.ID)
		return
	}
}
