.PHONY: run-bot
run-bot:
	go run ./cmd/bot

.PHONY: format
lint:
	golangci-lint run

.PHONY: format
format:
	gofumpt -l -w .