#!/bin/bash

# run migrations
./goose -dir ./db/sql sqlite3 ./betterbin.sqlite up
# run BetterBin
./BetterBin