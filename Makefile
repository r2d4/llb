# Variables
NAME := llb
ORG := r2d4
TAG := 1.0.3
REPO := ghcr.io
IMAGE := $(REPO)/$(ORG)/$(NAME):$(TAG)

GO_FILES := $(shell find . -name '*.go' -not -path "./vendor/*")
OUTPUT_DIR := _output
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

.PHONY: all
all: build

.PHONY: build
build:
	docker build -t $(IMAGE) .

.PHONY: bin
bin: $(OUTPUT_DIR)/go/bin/$(NAME)

.PHONY: push
push:
	docker push $(IMAGE)

$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)

$(OUTPUT_DIR)/go/bin/$(NAME): Dockerfile go.mod go.sum $(GO_FILES) $(OUTPUT_DIR)
	docker build -t $(IMAGE) --output type=local,dest=$(OUTPUT_DIR) . --platform $(GOOS)/$(GOARCH)

.PHONY: install
install: bin
	cp $(OUTPUT_DIR)/go/bin/$(NAME) $(HOME)/.go/bin
