.PHONY: all build run clean container

# Default target
all: build

build:
	@echo "Tidying up dependencies..."
	go mod tidy
	@echo "Building the binary..."
	go build -o bin/chitest

run: build
	@echo "Running the application..."
	bin/chitest

clean:
	@echo "Cleaning up..."
	rm -rf bin/

container: build
	@echo "Building the Docker container..."
	podman build -t docker.io/tranchida/chitest:latest -f Containerfile .