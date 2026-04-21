package db

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

var _ = Migrate(BookWork{})

// A literary work with all relevant metadata
type BookWork struct {
	// This struct intentionally does not inherit the BaseModel struct
	// The ID here is supplied by Google, not our own UUID
	gorm.Model
	ID string `gorm:"primaryKey"` // Google Books Volume ID

	Title         string        // E.g. A Woman Underground
	Subtitle      string        // E.g. A Cameron Winter Mystery
	Authors       SqlStringList `gorm:"type:text"`
	Publisher     string
	PublishedDate time.Time
	Version       string // As provided by Google Books

	Isbn13 string
	Isbn10 string

	Description string
	PageCount   int
	IsMature    bool
	Categories  SqlStringList `gorm:"type:text"`

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

// Gets all the editions matching this title and author
func (b *BookWork) Editions() ([]BookWork, error) {
	ret := []BookWork{}
	status := db.Where(&BookWork{
		Title:   b.Title,
		Authors: b.Authors,
	}).Find(&ret)
	return ret, status.Error
}

// Lists all copies from all editions
func (b *BookWork) AllCopies() (CopyList, error) {
	ret := []BookCopy{}
	status := db.Joins("book_copies").Where(&BookWork{
		Title:    b.Title,
		Subtitle: b.Subtitle,
		Authors:  b.Authors,
	}).Find(&ret)

	return ret, status.Error
}

// Strictly matches against this particular edition
func (b *BookWork) CopiesStrict() (CopyList, error) {
	ret := []BookCopy{}
	status := db.Where(&BookCopy{
		BookWorkID: b.ID,
	}).Find(&ret)
	return ret, status.Error
}

// TO-DO: Fix available copies

type CopyCount struct {
	Total     int
	Available int
}

type copyCountInter struct {
	Format BookFmtFlag
	Count  int
}

// Lists available copies
// If 'strict' is true, then only copies for this particular book ID will be returned
// If 'strict' is false, then copies matching the title and authors will be returned, too
func (b *BookWork) AvailableCopies(strict bool) (FormatsMap[CopyCount], error) {
	ids := []string{b.ID}
	if !strict {
		err := db.Where(&BookWork{
			Title:    b.Title,
			Subtitle: b.Subtitle,
			Authors:  b.Authors,
		}).Pluck("id", &ids).Error
		if err != nil {
			return nil, err
		}
	}

	var totalCounts []copyCountInter
	err := db.Model(&BookCopy{}).
		Where("status = ?", CopyStatusPublic).
		Where("book_work_id IN ?", ids).
		Group("format").
		Select("format, COUNT(*) as count").
		Scan(&totalCounts).Error
	if err != nil {
		return nil, err
	}

	checkedOutSubquery := db.Table("loans l").
		Select("1").
		Where("l.book_copy_id = book_copies.id").
		Where("l.date_returned = ?", NilTime)

	var availableCounts []copyCountInter
	err = db.Model(&BookCopy{}).
		Where("status = ?", CopyStatusPublic).
		Where("book_work_id IN ?", ids).
		Where("NOT EXISTS (?)", checkedOutSubquery).
		Group("format").
		Select("format, COUNT(*) as count").
		Scan(&availableCounts).Error
	if err != nil {
		return nil, err
	}

	ret := FormatsMap[CopyCount]{}
	for _, tc := range totalCounts {
		ret[tc.Format] = CopyCount{Total: tc.Count}
	}
	for _, ac := range availableCounts {
		if c, ok := ret[ac.Format]; ok {
			c.Available = ac.Count
			ret[ac.Format] = c
		} else {
			ret[ac.Format] = CopyCount{
				Total:     ac.Count,
				Available: ac.Count,
			}
		}
	}

	return ret, nil
}
