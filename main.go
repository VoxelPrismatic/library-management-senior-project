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
	if code := goRunCommand("templ", "generate", "-path", "."); code != 0 {
		os.Exit(code)
	}

	args := []string{"run", "main.go", "-templ-done"}
	args = append(args, os.Args[1:]...)
	fmt.Println(args)
	if code := goRunCommand("go", args...); code != 0 {
		os.Exit(code)
	}
}

func goRunCommand(bin string, args ...string) int {
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			return exitErr.ExitCode()
		} else {
			panic("unreachable")
		}
	}
	return 0
}
