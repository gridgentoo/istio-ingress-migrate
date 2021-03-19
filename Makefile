.ONESHELL:
SHELL := /bin/bash
GOBIN ?= $(GOPATH)/bin
MODULE = github.com/istio-ecosystem/istio-ingress-migrate
HUB ?= gcr.io/istio-release
VERSION ?= 0.0.1
export GO111MODULE ?= on

.PHONY: format
format: $(GOBIN)/goimports
	@go mod tidy
	@goimports -l -w -local $(MODULE) .

.PHONY: install
install:
	@go install

.PHONY: docker
docker:
	docker buildx build . -t ${HUB}/istio-ingress-migrate -t ${HUB}/istio-ingress-migrate:${VERSION} --load

.PHONY: push
push:
	docker buildx build . -t ${HUB}/istio-ingress-migrate -t ${HUB}/istio-ingress-migrate:${VERSION} --push
