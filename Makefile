## test: runs all tests
test:
	@go test -v ./...

## cover: opens coverage in browser
cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

## coverage: displays test coverage
coverage:
	@go test -cover ./...

## build_cli: builds the jazz command line tool and copies it to myapp
build_cli:
	@go build -o ../myapp/jazz ./cmd/cli

## build: builds the command line tool to dist directory
build:
	@go build -o ./dist/jazz ./cmd/cli