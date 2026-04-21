package fetch

import (
	"slices"
	"time"
	"voxelprismatic/library-management-senior-project/db"
)

// https://developers.google.com/books/docs/v1/using

type GBooksVolDetails struct {
	Kind       string           `json:"kind"`
	ID         string           `json:"id"`
	Etag       string           `json:"etag"`
	SelfLink   string           `json:"selfLink"`
	VolumeInfo GBooksVolInfo    `json:"volumeInfo"`
	SaleInfo   GBooksSaleInfo   `json:"saleInfo"`
	AccessInfo GBooksAccessInfo `json:"accessInfo"`
}

type GBooksVolInfo struct {
	Title               string              `json:"title"`
	Subtitle            string              `json:"subtitle"`
	Authors             []string            `json:"authors"`
	Publisher           string              `json:"publisher"`
	PublishedDate       string              `json:"publishedDate"`
	Description         string              `json:"description"`
	IndustryIdentifiers []GBooksIndustryIDs `json:"industryIdentifiers"`
	ReadingModes        GBooksReadingModes  `json:"readingModes"`
	PageCount           int                 `json:"pageCount"`
	PrintedPageCount    int                 `json:"printedPageCount"`
	MaturityRating      string              `json:"maturityRating"`
	Dimensions          GBooksDimensions    `json:"dimensions"`
	PrintType           string              `json:"printType"`
	MainCategory        string              `json:"mainCategory"`
	Categories          []string            `json:"categories"`
	AverageRating       float32             `json:"averageRating"`
	RatingsCount        int                 `json:"ratingsCount"`
	ContentVersion      string              `json:"contentVersion"`
	ImageLinks          GBooksImgLinks      `json:"imageLinks"`
	Language            string              `json:"language"`
	InfoLink            string              `json:"infoLink"`
	PreviewLink         string              `json:"previewLink"`
	CanonicalVolumeLink string              `json:"canonicalVolumeLink"`
}

type GBooksReadingModes struct {
	Text  bool `json:"text"`
	Image bool `json:"image"`
}

type GBooksIndustryIDs struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

type GBooksDimensions struct {
	Height    string `json:"height"`
	Width     string `json:"width"`
	Thickness string `json:"thickness"`
}

type GBooksImgLinks struct {
	SmallThumbnail string `json:"smallThumbnail"`
	Thumbnail      string `json:"thumbnail"`
	Small          string `json:"small"`
	Medium         string `json:"medium"`
	Large          string `json:"large"`
	ExtraLarge     string `json:"extraLarge"`
}

type GBooksSaleInfo struct {
	Country     string          `json:"country"`
	Saleability string          `json:"saleability"`
	IsEbook     bool            `json:"isEbook"`
	ListPrice   GBooksPriceInfo `json:"listPrice"`
	RetailPrice GBooksPriceInfo `json:"retailPrice"`
	BuyLink     string          `json:"buyLink"`
}

type GBooksPriceInfo struct {
	Amount       float32 `json:"amount"`
	CurrencyCode string  `json:"currencyCode"`
}

type GBooksAccessInfo struct {
	Country                string           `json:"country"`
	Viewability            string           `json:"viewability"`
	Embeddable             bool             `json:"embeddable"`
	PublicDomain           bool             `json:"publicDomain"`
	TextToSpeechPermission string           `json:"textToSpeechPermission"`
	Epub                   GBooksFormatInfo `json:"epub"`
	Pdf                    GBooksFormatInfo `json:"pdf"`
	AccessViewStatus       string           `json:"accessViewStatus"`
}

type GBooksFormatInfo struct {
	IsAvailable  bool   `json:"IsAvailable"`
	AcsTokenLink string `json:"acsTokenLink"`
}

type GBooksVolSearch struct {
	Kind       string             `json:"kind"`
	TotalItems int                `json:"totalItems"`
	Items      []GBooksVolDetails `json:"items"`
}

func (b GBooksVolDetails) ToLocalStruct() db.BookWork {
	v := b.VolumeInfo
	pubDate, _ := time.Parse("2006-01-02", v.PublishedDate)

	var isbn10, isbn13 string
	for _, id := range v.IndustryIdentifiers {
		switch id.Type {
		case "ISBN_10":
			isbn10 = id.Identifier
		case "ISBN_13":
			isbn13 = id.Identifier
		}
	}

	categories := db.SqlStringList(v.Categories[:])
	if v.MainCategory != "" && !slices.Contains(v.Categories, v.MainCategory) {
		categories = append(categories, v.MainCategory)
	}

	return db.BookWork{
		ID:            b.ID,
		Title:         v.Title,
		Subtitle:      v.Subtitle,
		Authors:       db.SqlStringList(v.Authors),
		Publisher:     v.Publisher,
		PublishedDate: pubDate,
		Version:       v.ContentVersion,
		Isbn13:        isbn13,
		Isbn10:        isbn10,
		Description:   v.Description,
		PageCount:     max(v.PageCount, v.PrintedPageCount),
		IsMature:      v.MaturityRating == "MATURE",
		Categories:    categories,
		CoverThumb:    v.ImageLinks.Thumbnail,
		CoverImage:    v.ImageLinks.Large,
	}
}
