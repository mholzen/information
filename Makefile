test:
	(cd triples; go test -v)

run:
	go run server/server.go
	open http://localhost:1323/html/examples.json