BUILDFLAGS = -buildmode=pie

# update go dependencies
update:
	go get -u ./cmd
	go mod tidy

mock:
	@mockery --all

# run linter
lint:
	golangci-lint run ./...

# run linter and fix issues if possible
lintfix:
	golangci-lint run --fix ./...

# run unit tests
test:
	@go test ${BUILDFLAGS} -coverprofile=cover.out ./...
	@go tool cover -func=cover.out
	-@rm -f cover.out

# run ttm, note: make doesn't understand exit code 130 and sets it == 1
run:
	@go run ${BUILDFLAGS} ./cmd || exit 0

install: build
	@mv ./ttm ${HOME}/go/bin/
	@echo "ttm has been installed."

# build ttm
build:
	go build ${BUILDFLAGS} -v -o ttm ./cmd
