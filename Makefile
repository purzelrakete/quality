.PHONY: all build test deps install clean

PKGS := \
github.com/purzelrakete/quality \
github.com/purzelrakete/quality/crawl

all: deps build test install

build:
	go build -v $(PKGS)

test:
	go test -v

deps:
	go get -v $(PKGS)

install:
	go install -v $(PKGS)

clean:
	go clean $(PKGS)
