GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=distributed-consensus

build: 
    $(GOBUILD) -o $(BINARY_NAME) -v
test: 
    $(GOTEST) -v ./...

run:
    $(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)