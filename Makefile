.PHONY: all test lint coverage integration env

COVER_PROFILE=cover.out

all: test lint

test:
	go test -race -cover

lint:
	golint

	@[ ! -f ./bin/golangci-lint ] && curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh \
		| sh -s -- -b ./bin v1.21.0 || true
	./bin/golangci-lint run

coverage:
	go test -coverprofile $(COVER_PROFILE) && go tool cover -html=$(COVER_PROFILE)
