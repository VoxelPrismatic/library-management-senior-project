package main

import (
	"flag"
	"fmt"
	"voxelprismatic/library-management-senior-project/db"
	"voxelprismatic/library-management-senior-project/fetch"
)

func main() {
	gapiToken := flag.String("gapi-token", "", "Google API token")
	flag.Parse()

	fetch.SetAPIToken(*gapiToken)

	data, err := fetch.GBooksVolume("Bj6VEAAAQBAJ")
	if err != nil {
		panic(err)
	}

	fmt.Println(data)
	book := data.ToLocalStruct()
	db.Save(&book)
}
