cwd = $(shell pwd)

run:
	env ROOT=$(cwd) go run server/server.go

open:
	open http://localhost:1323/html/data/examples.jsonc

test:
	go test ./... | grcat ~/.grc/go.conf
