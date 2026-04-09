package scrape

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
	Authors             []string            `json:"authors"`
	Publisher           string              `json:"publisher"`
	PublishedDate       string              `json:"publishedDate"`
	Description         string              `json:"description"`
	IndustryIdentifiers []GBooksIndustryIDs `json:"industryIdentifiers"`
	PageCount           int                 `json:"pageCount"`
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
	CanonicalVolumeLink string              `json:"canonicalVolumeLink"`
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
	Kind  string          `json:"kind"`
	Items []GBooksVolInfo `json:"items"`
}
