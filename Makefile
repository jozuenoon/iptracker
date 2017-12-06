SHELL := /bin/bash

build:
	go build -o bin/iptracker

install: build
	$(shell cp bin/iptracker $$HOME/bin)
	