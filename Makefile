cwd = $(shell pwd)

run:
	env ROOT=$(cwd) go run server/server.go

open:
	open http://localhost:1323/

test:
	env ROOT=$(cwd) go test ./... | grcat ~/.grc/go.conf
