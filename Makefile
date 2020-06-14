.PHONY: help clean build debug install test tidy sniff .travis

help:                   # Displays this list
	@echo; grep "^[a-z][a-zA-Z0-9_<> -]\+:" Makefile | sed -E "s/:[^#]*?#?(.*)?/\r\t\t\1/" | sed "s/^/ make /"; echo

clean:                  # Removes build/test artifacts
	@find . -type f | grep "\.out$$"  | xargs -I{} rm {};
	@find . -type f | grep "\.html$$" | xargs -I{} rm {};
	@find . -type f | grep "\.test$$" | xargs -I{} rm {};

build: clean            # Builds a static binary to ./bin/typex
	@go build -trimpath -ldflags "-s -w" -o ./bin/typex .

debug: clean            # Starts debugger [:2345] with ./bin/typex
	@go build -gcflags "all=-N -l" -o ./bin/typex .
	dlv --listen=:2345 --headless=true --api-version=2 exec ./bin/typex $(ARGS)

install: clean          # Compiles and installs typex in Go environment
	@go install -trimpath -ldflags "-s -w" .

test: clean             # Runs tests, reports coverage
	@go test -v -count=1 -covermode=atomic -coverprofile=./coverage.out -coverpkg=./internal/... ./internal/...
	@go tool cover -html=./coverage.out -o ./coverage.html && echo "coverage: <file://$(PWD)/coverage.html>"

tidy:                  # Formats source files, cleans go.mod
	@gofmt -w .
	@go mod tidy

sniff:                  # Checks format and runs linter (void on success)
	@gofmt -d .
	@2>/dev/null revive -config .revive.toml ./... || echo "get a linter first:  go install github.com/mgechev/revive"

.travis:                # Travis CI (see .travis.yml), runs tests
ifndef TRAVIS
	@echo "Fail: requires Travis runtime"
else
	@$(MAKE) test --no-print-directory && goveralls -coverprofile=./coverage.out -service=travis-ci
endif
