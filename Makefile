
SHELL := /bin/bash

# The name of the executable (default is current directory name)
TARGET := pldrctl
.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
VERSION := 0.0.1
BUILD := `git rev-parse HEAD`

# Operating System Default (LINUX)
TARGETOS=linux

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -s"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

DOCKERTAG=latest

.PHONY: all build clean install uninstall fmt simplify check run

all: check install

$(TARGET): $(SRC)
	@go build $(LDFLAGS) -o $(TARGET)

build: $(TARGET)
	@true

clean:
	@rm -f $(TARGET)

install:
	@echo Building and Installing project
	@go install $(LDFLAGS)

uninstall: clean
	@rm -f $$(which ${TARGET})

fmt:
	@gofmt -l -w $(SRC)

docker:
	@GOOS=$(TARGETOS) make build
	@mv $(TARGET) ./dockerfile
	@docker build -t $(TARGET):$(DOCKERTAG) ./dockerfile/
	@rm ./dockerfile/$(TARGET)
	@echo New Docker image created

oui:
	@echo Pulling latest updates from the oui database
	@mkdir -p pkg/oui
	@echo -e "package oui\n\nconst ouiLookup=\``curl http://standards-oui.ieee.org/oui.txt | grep "(hex)" | sort |  cut -c 1-8,17- | tr -d '\`'`\n\`" > pkg/oui/lookup.go

all_releases:
	@make release_darwin
	@make release_linux
	@make release_linux_arm
	@make release_windows

release_darwin:
	@echo Creating Darwin Build
	@GOOS=darwin make build
	@zip -9 -r $(TARGET)-darwin-$(VERSION).zip ./$(TARGET)
	@rm $(TARGET)

release_linux:
	@echo Creating Linux amd64 Build
	@GOOS=linux GOARCH=amd64 make build
	@zip -9 -r $(TARGET)-linux-amd64-$(VERSION).zip ./$(TARGET) 
	@rm $(TARGET)

release_linux_arm:
	@echo Creating Linux Arm Build
	@GOOS=linux GOARCH=arm64 make build
	@zip -9 -r $(TARGET)-linux-armv64-$(VERSION).zip ./$(TARGET) 
	@rm $(TARGET)

release_windows:
	@echo Creating Windows Build
	@GOOS=windows make build
	@zip -9 -r $(TARGET)-win64-$(VERSION).zip ./$(TARGET) 
	@rm $(TARGET)

simplify:
	@gofmt -s -l -w $(SRC)

check:
	@test -z $(shell gofmt -l main.go | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done
	@go tool vet ${SRC}

run: install
	@$(TARGET)
