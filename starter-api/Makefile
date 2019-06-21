# Copyright 2018 Payfazz Agent Authors.

# The default docker image name
IMAGE := payfazz-cashout-api
 
# The docker registry ip
REGISTRY := docker.fazzfinancial.com/payfazz

# This version-strategy uses git tags to set the version string
VERSION := $(shell git describe --tags --always --dirty)

migrate:
	@go run -v ./databases/migrate.go

test:
	@go test -v -cover -p 1 ./...

container: Dockerfile
	@echo "building the container..."
	@docker build --label "version=$(VERSION)" -t $(IMAGE):$(VERSION) .
 
# before the push, you need to login to docker.fazzfinancial.com: 
push-docker: .push-docker
.push-docker:
	docker tag $(IMAGE):$(VERSION) $(REGISTRY)/$(IMAGE):$(VERSION)
	docker push $(REGISTRY)/$(IMAGE):$(VERSION)
	echo "pushed: $(REGISTRY)/$(IMAGE):$(VERSION)"

clean:
	rm -rf .container-* .dockerfile-* .push-*
	rm -rf .go bin
