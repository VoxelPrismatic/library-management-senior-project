package main

import (
	"fmt"
	"voxelprismatic/library-management-senior-project/fetch"
)

func main() {
	data, err := fetch.GBooksVolume("Bj6VEAAAQBAJ")
	if err != nil {
		panic(err)
	}

	fmt.Println(data)
}
