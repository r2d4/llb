.PHONY: all
all: build

.PHONY: build
build:
	cd build && tsx src/main
