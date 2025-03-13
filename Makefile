TARGET := bin
INSTALL_DIR := /usr/local/bin
BINARY := bitsmuggler

all: build

build:
	@echo "Building bitsmuggler..."
	go build -o $(TARGET)/$(BINARY)
	@echo "Built!"

install: build
	install -m 0755 $(TARGET)/$(BINARY) $(INSTALL_DIR)
	@echo "Done!"

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY)
	@echo "Uninstalled!"

clean:
	rm -rdf bin/
	@echo "Done!"

test:
	@echo "Running tests..."
	@go test -v ./yts/ ./opensubs/

format:
	@echo "Formatting..."
	@gofmt -w .
	@echo "Done!"

.PHONY: all build install uninstall clean
