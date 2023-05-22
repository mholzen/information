cwd = $(shell pwd)

run:
	env ROOT=$(cwd) go run server/server.go
	open http://localhost:1323/html/examples.json

test:
	(cd triples; make test)
