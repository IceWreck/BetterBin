#!/bin/bash

# run migrations
touch /home/betterbin/app/data/betterbin.sqlite
./goose -dir ./db/sql sqlite3 /home/betterbin/app/data/betterbin.sqlite up
# run BetterBin
./BetterBin -p 8963 -d /home/betterbin/app/data/betterbin.sqlite