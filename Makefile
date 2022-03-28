
## testing

.PHONY: coverage.out
coverage.out:
	go test -race -covermode=atomic -coverprofile=$@ ./...

coverage.html: coverage.out
	go tool cover -html $< -o $@

.PHONY: test
test: coverage.out ## runs tests
