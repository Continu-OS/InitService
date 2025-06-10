# Pfade
SRC_DIR=src/syscored/proc
OUTPUT_DIR=output

# Standard-Build-Ziel
all: build

build:
	@echo "🔨 Building syscored from src/proc/syscored..."
	GOOS=linux GOARCH=arm64 go build -o $(OUTPUT_DIR)/syscored ./src/proc/syscored

clean:
	@echo "🧹 Cleaning up..."
	rm -rf $(OUTPUT_DIR)/syscored

run: build
	@echo "🚀 Running syscored..."
	./$(OUTPUT_DIR)/syscored

.PHONY: all build clean run
