# Makefile im Hauptverzeichnis

# Binary-Name ohne Pfad
BINARY_NAME=initservice

# Pfade
SRC_DIR=src/init/proc
OUTPUT_DIR=output

# Standard-Build-Ziel
all: build

build:
	@echo "🔨 Building $(BINARY_NAME) from $(SRC_DIR)..."
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME) ./$(SRC_DIR)

clean:
	@echo "🧹 Cleaning up..."
	rm -rf $(OUTPUT_DIR)/$(BINARY_NAME)

run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	./$(OUTPUT_DIR)/$(BINARY_NAME)

.PHONY: all build clean run
