cwd = $(shell pwd)

define run
    @env ROOT=$(cwd) go run server/server.go
endef


run:
	$(call run)

start:
	@echo "Starting server..."
	env ROOT=$(cwd) sh -c 'go run server/server.go & echo "$$!" > .pid'

stop:
	kill -term $(shell cat .pid)

watch:
	@fswatch -o . | xargs -n1 -I{} make stop && make start

open:
	open http://localhost:1323/

test:
	env ROOT=$(cwd) go test ./... | grcat ~/.grc/go.conf
