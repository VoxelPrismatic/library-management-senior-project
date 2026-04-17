package book

import (
	"fmt"
	"strings"
	"time"

	"voxelprismatic/library-management-senior-project/db"

	"gorm.io/gorm"
)

var _ = db.Migrate(BookWork{})

// A literary work with all relevant metadata
type BookWork struct {
	gorm.Model
	ID string `gorm:"primaryKey"` // Google Books Volume ID

	Title         string           // E.g. A Woman Underground
	Subtitle      string           // E.g. A Cameron Winter Mystery
	Authors       db.SqlStringList `gorm:"type:text"`
	Publisher     string
	PublishedDate time.Time
	Version       string // As provided by Google Books

	Isbn13 string
	Isbn10 string

	Description string
	PageCount   int
	IsMature    bool
	Categories  db.SqlStringList `gorm:"type:text"`

	CoverThumb string
	CoverImage string
}

func (b *BookWork) Tags() []string {
	set := map[string]bool{}
	for _, category := range b.Categories {
		for subcategory := range strings.SplitSeq(category, "/") {
			fmt.Println(subcategory)
			subcategory = strings.TrimSpace(subcategory)
			set[subcategory] = true
		}
	}

	ret := make([]string, len(set))
	i := 0
	for k := range set {
		ret[i] = k
		i++
	}
	return ret
}

type BookVariants map[string][]BookWork

func (v *BookVariants) Add(b BookWork) {
	id := b.Isbn13
	arr, ok := (*v)[id]
	if !ok {
		id = b.Isbn10
		arr, ok = (*v)[id]
	}
	if ok {
		for _, e := range arr {
			if e.ID == b.ID {
				// Duplicate in search for whatever reason
				return
			}
		}
		(*v)[id] = append(arr, b)
		return
	}

	id = b.Isbn13
	if id == "" {
		id = b.Isbn10
	}
	if id == "" {
		id = b.ID
	}
	(*v)[id] = []BookWork{b}
}
