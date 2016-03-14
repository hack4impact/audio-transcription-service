NON_VENDOR := $(shell go list ./... | grep -v /vendor/)

all: test

test:
	mocktest -v $(NON_VENDOR)
