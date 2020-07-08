PROJECTNAME="url-shortener"
LDFLAGS=-ldflags "-s -w"

clean:
	rm -rf .serverless ./bin

run:
	go run main.go

dev:
	env GOOS=linux go build -o $(PROJECTNAME) main.go

build:
	env GOOS=linux go build $(LDFLAGS) -o bin/$(PROJECTNAME) main.go

test:
	go test -v ./...

coverage:
	go test -v ./... -coverprofile cover.out && go tool cover -html=cover.out

deploy: clean build
	./node_modules/.bin/sls deploy --verbose
