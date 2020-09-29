SHELL := /bin/bash
.PHONY: build
build:
	@echo 'Building Docker image'
	docker build -t crud-using-chi .

run:
	@echo 'Running application'
	docker-compose up -d