run:
	@echo "Starting up..."
	@trap 'echo "Cleaning up..."; rm -rf data-dir-*' EXIT
	@go run cmd/app/main.go
