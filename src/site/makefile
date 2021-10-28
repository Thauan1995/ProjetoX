#!/usr/bin/env make -f

appid=estudos-312813

deploy-prod: build-server
	go test ./...
	gcloud app deploy app.yaml \
	--project=${appid} \

deploy-branch: build-server
	gcloud app deploy app.yaml \
	--project=${appid} \
	--version $$(git rev-parse --abbrev-ref HEAD | tr '[:upper:]' '[:lower:]') \
	--no-promote --quiet

build-server:
	go vet ./... && \
	go build -v ./...