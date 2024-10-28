.PHONY: build clean run

# Binary name
BINARY_NAME=distributor-cli

# Build directory
BUILD_DIR=build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Main build target
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/distributor-cli

# Clean build files
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run the application (example usage)
run: build
	./$(BUILD_DIR)/$(BINARY_NAME) \
		-csv ./data/cities.csv \
		-perm ./data/permissions.txt \
		-dist "DISTRIBUTOR2" \
		-check "RANCH-JH-IN"
