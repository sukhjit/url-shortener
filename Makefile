DIST_DIR=$(shell pwd)/dist
BIN_DIR=$(DIST_DIR)/bin
TEST_DIR=$(DIST_DIR)/tests
GO_BUILD=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w"
LINT_VERSION=v1.50.0
LINT_URL=https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh
LINTER=$(BIN_DIR)/golangci-lint

clean:
	rm -rf .serverless ./bin

dev:
	@cicd/bin/local-dev.sh

lint-install:
	@mkdir -p $(BIN_DIR)
	@[ -f "$(LINTER)" ] || curl -sSfL $(LINT_URL) | sh -s -- -b $(BIN_DIR) $(LINT_VERSION)

lint: lint-install
	@$(LINTER) run -v ./...

build:
	$(GO_BUILD) -o bin/url-shortener main.go

test:
	@mkdir -p $(TEST_DIR)
	@go clean -testcache
	@go test \
		-coverpkg=./... \
		-coverprofile=c.out \
		-outputdir=$(TEST_DIR) \
		-race \
		-failfast \
		./...
	@go tool cover -html=$(TEST_DIR)/c.out -o $(TEST_DIR)/c.html

deploy: clean build
	npm install
	./node_modules/.bin/sls deploy --verbose
