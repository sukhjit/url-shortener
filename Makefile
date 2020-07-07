PROJECTNAME="url-shortener"
LDFLAGS=-ldflags "-s -w"

run:
	go run main.go

dev:
	env GOOS=linux go build -o $(PROJECTNAME) main.go

build:
	env GOOS=linux go build $(LDFLAGS) -o $(PROJECTNAME) main.go

test:
	go test -v ./...

coverage:
	go test -v ./... -coverprofile cover.out && go tool cover -html=cover.out
