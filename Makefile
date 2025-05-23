CLIENT_DIR = client
DIST_DIR = $(CLIENT_DIR)/dist
GO_BINARY = server
OUTPUT_DIR = build

.PHONY: all
all: build

.PHONY: build
build: build-client build-server

.PHONY: build-client
build-client:
	@echo "Building client..."
	cd $(CLIENT_DIR) && npm ci && npm run build

.PHONY: build-server
build-server:
	@echo "Building server..."
	mkdir -p $(OUTPUT_DIR)
	go build -o $(OUTPUT_DIR)/$(GO_BINARY)

.PHONY: clean
clean:
	rm -rf $(DIST_DIR) $(OUTPUT_DIR)

.PHONY: run
run: build
	./$(OUTPUT_DIR)/$(GO_BINARY)

