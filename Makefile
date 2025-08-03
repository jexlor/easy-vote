# Makefile

MAIN := ./cmd/main.go
GO := go

.PHONY: run

run:
	$(GO) run $(MAIN)
