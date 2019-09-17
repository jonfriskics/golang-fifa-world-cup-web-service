GOCMD=go
GOTEST=$(GOCMD) test
GORUN=${GOCMD} run

all: test

start: 
				${GORUN} server.go
test:
				$(GOTEST) -v ./handlers/*
