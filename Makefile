 # Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
RICE=rice embed-go -i ./wsgatherer
BINARY_NAME=main
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=cmd/main.go

build: 
	$(RICE)
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
# deps: