package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GBooksSearch(search string) (*GBooksVolSearch, error) {
	resp, err := http.Get("https://www.googleapis.com/books/v1/volumes?q=" + search)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := &GBooksVolSearch{}
	err = json.Unmarshal(body, ret)
	return ret, err
}

func GBooksIsbnLookup(isbn int) (*GBooksVolSearch, error) {
	return GBooksSearch(fmt.Sprintf("isbn:%d", isbn))
}

func GBooksVolume(volume string) (*GBooksVolDetails, error) {
	resp, err := http.Get("https://www.googleapis.com/books/v1/volumes/" + volume + "?projection=full")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))

	ret := &GBooksVolDetails{}
	err = json.Unmarshal(body, ret)
	return ret, err
}
