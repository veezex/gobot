.PHONY: run-bot
run-bot:
	go run ./cmd/bot

.PHONY: build-bot
build-bot:
	go build ./cmd/bot

.PHONY: format
lint:
	golangci-lint run

.PHONY: format
format:
	gofumpt -l -w .

.PHONY: update
update:
	go get -u ./...
	go mod tidy