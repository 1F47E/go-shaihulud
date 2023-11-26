.PHONY: server client

server:
	@echo "Starting server..."
	@go run cmd/app/main.go srv

client:
	@echo "Starting client..."
	@go run cmd/app/main.go cli

clean:
	@echo "Cleaning up..."
	@sudo rm -rf data-dir-*
