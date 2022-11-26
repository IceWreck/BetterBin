# run code in dev mode
run:
	go run .

# generate a binary
build:
	CGO_ENABLED=0 go build -o out/BetterBin .

# build a container image
build-container:
	podman build -t icewreck/betterbin .

# run created container (without attaching volume)
run-container:
	podman run -p 8080:8963 --rm -ti icewreck/betterbin

# required when you need to build with older glibc (for older servers)
start-old-container:
	podman run -it --rm -v $(shell pwd):/src:z docker.io/library/golang:1.15

# USAGE: make migrate_new MNAME=whatever_name
migrate_new:
	goose -dir ./db/sql sqlite3 ./betterbin.sqlite create $(MNAME) sql

# up all migrations
migrate_up:
	goose -dir ./db/sql sqlite3 ./betterbin.sqlite up

# down one migration
migrate_down:
	goose -dir ./db/sql sqlite3 ./betterbin.sqlite down

# migration status
migrate_status:
	goose -dir ./db/sql sqlite3 ./betterbin.sqlite status
