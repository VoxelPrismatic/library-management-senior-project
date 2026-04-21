.PHONY: templ test-db test-fetch
run: templ
	go run main.go

templ:
	templ generate -path .

test: templ test-db test-fetch

test-db:
	go test ./db

test-fetch:
	go test ./fetch
