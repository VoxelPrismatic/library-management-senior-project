.PHONY: run templ
run: templ
	go run main.go

templ:
	templ generate -path .
