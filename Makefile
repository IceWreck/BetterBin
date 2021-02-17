# run code in dev mode
run:
	go run .

# generate a binary
build:
	go build -o out/BetterBin .

# required to build with older glibc
container:
	podman run -it --rm -v $(shell pwd):/src:z docker.io/library/golang:1.14.15-stretch

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