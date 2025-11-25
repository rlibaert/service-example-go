all: lint test dist

.PHONY: lint
lint:
	golangci-lint run -v

.PHONY: test
test:
	go test ./... -cover
	go test ./... -race

cover_%.html: cover_%.out
	go tool cover -html=$< -o $@

cover_%.out: %/*
	go test ./$* -covermode=count -coverprofile=$@

dist: vendor
	docker run --rm \
	-w /workdir \
	-v ${PWD}:/workdir \
	-v /var/run/docker.sock:/var/run/docker.sock \
	goreleaser/goreleaser build --clean --snapshot

vendor:
	go mod vendor

.PHONY: clean
clean:
	rm -rf vendor dist
