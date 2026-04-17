package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var token string

/*
Search for books. Returns a "list" (google's search object is a list inside an object)

Special search tags:
 1. `intitle:` Return results where the text following this word is found in the title.
 2. `inauthor:` …found in the author's name.
 3. `inpublisher:` …found in the publisher's name.
 4. `subject:` …found in the list of categories.
 5. `isbn:` Where the immediate next word matches the ISBN.
 6. `lccn:` …matches the Library of Congress Control number,
 7. `oclc:` …matches the Online Computer Library Center number.
*/
func GBooksSearch(search string) (*GBooksVolSearch, error) {
	uri := "https://www.googleapis.com/books/v1/volumes?q=" + url.QueryEscape(search)
	if token != "" {
		uri += "&key=" + token
	}
	resp, err := http.Get(uri)
	fmt.Println(uri)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := &GBooksVolSearch{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		fmt.Printf("\x1b[91;1m%s\x1b[0m\n", uri)
		fmt.Println(string(body))
	}
	return ret, err
}

// Google Books Volume ID
func GBooksVolume(volume string) (*GBooksVolDetails, error) {
	uri := "https://www.googleapis.com/books/v1/volumes/" + volume + "?projection=full"
	if token != "" {
		uri += "&key=" + token
	}
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := &GBooksVolDetails{}
	err = json.Unmarshal(body, ret)
	return ret, err
}

func SetAPIToken(newToken string) {
	token = newToken
	fmt.Println("set token")
}
