all: lint test dist

.PHONY: generate
generate:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/proto/service.proto

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
