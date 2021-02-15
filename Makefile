run:
	go run .

build:
	go build -o dump/BetterBin .

# required to build with older glibc
container:
	podman run -it --rm -v $(shell pwd):/src:z docker.io/library/golang:1.14.15-stretch