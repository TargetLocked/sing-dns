fmt:
	@gofumpt -l -w .
	@gofmt -s -w .
	@gci write --custom-order -s standard -s "prefix(github.com/sagernet/)" -s "default" .

fmt_install:
	go install -v mvdan.cc/gofumpt@latest
	go install -v github.com/daixiang0/gci@latest

lint:
	GOOS=linux golangci-lint run ./...
	GOOS=android golangci-lint run ./...
	GOOS=windows golangci-lint run ./...
	GOOS=darwin golangci-lint run ./...
	GOOS=freebsd golangci-lint run ./...

lint_install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

test:
	go test -v ./...

diff:
	git fetch upstream
	git diff origin/dev..upstream/dev

upgrade:
	git fetch upstream
	git branch -f dev upstream/dev
	git rebase origin/dev dance-crate --onto upstream/dev

push:
	git push -f origin dev
	git push -f origin dance-crate

COMMIT_INFO = $(shell git log -1 --format="%H %ct")
MODVER_HASH = $(shell echo $(word 1, $(COMMIT_INFO)) | cut -c 1-12)
MODVER_TIME = $(shell date -u -d @$(word 2,$(COMMIT_INFO)) +"%Y%m%d%H%M%S")
MODVER = v0.0.0-$(MODVER_TIME)-$(MODVER_HASH)

modver:
	@echo $(MODVER)