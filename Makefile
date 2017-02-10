GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test

.PHONY: testdata

testdata:
	$(GOBUILD) -v ./
	./goflow -dir=./testdata -out=./testdata