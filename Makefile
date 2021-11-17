PROJECTNAME=url-shortener
GO=go
GO_TEST=$(GO) test -race -failfast ./handler ./model ./repo
GO_COVER=$(GO) tool cover
GO_BUILD=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -trimpath -ldflags "-s -w"

clean:
	rm -rf .serverless ./bin

dev:
	@cicd/bin/local-dev.sh

lint:
	@cicd/bin/lint.sh

build:
	$(GO_BUILD) -o bin/$(PROJECTNAME) main.go

test:
	@$(GO_TEST)

coverage:
	@$(GO_TEST) -coverprofile=c.out -coverpkg=./...
	@$(GO_COVER) -func=c.out
	@$(GO_COVER) -html=c.out -o coverage.html
	@echo "\nWriting coverage to coverage.html file\n"

deploy: clean build
	./node_modules/.bin/sls deploy --verbose
