#!/usr/bin/env bash

.PHONY: all latest

default: all

all:
	go run -v web.go

test:
	go test -v ./...

latest:
	curl -o latest.dump `heroku pg:backups public-url --app brt-backend`
	pg_restore --verbose --clean --no-acl --no-owner -d brt latest.dump
