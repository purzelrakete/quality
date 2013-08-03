.PHONY: all build test deps install clean

PKGS := \
github.com/purzelrakete/quality \
github.com/purzelrakete/quality/crawl

all: deps build test install

build:
	go build -v $(PKGS)

test:
	go test -v

coverage:
	goveralls -service drone.io $$COVERALLS_TOKEN

deps:
	go get -v $(PKGS)
	go get github.com/axw/gocov/gocov
	go get github.com/mattn/goveralls

install:
	go install -v $(PKGS)

clean:
	go clean $(PKGS)
