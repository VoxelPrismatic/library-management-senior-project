package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"voxelprismatic/library-management-senior-project/fetch"
	"voxelprismatic/library-management-senior-project/router"
)

func main() {
	gapiToken := flag.String("gapi-token", "", "Google API token")
	listenAddr := flag.String("addr", "0.0.0.0", "Listen address")
	listenPort := flag.Int("port", 3000, "Listen port")
	didTempl := flag.Bool("templ-done", false, "Do not regenerate templ before launching")

	flag.Parse()

	if !*didTempl {
		goRebuildUrself()
		return
	}

	fetch.SetAPIToken(*gapiToken)

	fmt.Println("Starting server...")
	http.HandleFunc("/", router.Router)
	listen := fmt.Sprintf("%s:%d", *listenAddr, *listenPort)
	fmt.Printf("Listening on %s\n", listen)
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		panic(err)
	}

	/*
		data, err := fetch.GBooksVolume("Bj6VEAAAQBAJ")
		if err != nil {
			panic(err)
		}

		fmt.Println(data)
		book := data.ToLocalStruct()
		db.Save(&book)
	*/
}

func goRebuildUrself() {
	cmd := exec.Command("templ", "generate", "-path", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			os.Exit(exitErr.ExitCode())
		} else {
			panic("unreachable")
		}
	}
	args := []string{"run", "main.go", "-templ"}
	args = append(args, os.Args...)
	cmd = exec.Command("go", args...)
	_ = cmd.Run()
}
