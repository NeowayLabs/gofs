test:
	go test -race -coverprofile=coverage.txt -covermode=atomic

deps:
	go get -t -v ./...
