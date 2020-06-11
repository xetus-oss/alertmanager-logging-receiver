VERSION ?= latest
REPOSITORY_NAME = xetusoss/alertmanager-logging-receiver
PUSH_TAG := ${REPOSITORY_NAME}:${VERSION}
IMAGE_TAG := alertmanager-logging-receiver:${VERSION}

.PHONY: help
help: Makefile
	@echo
	@echo " Available targets"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'

.PHONY: clean
## clean: remove build outputs
clean:
	rm -rf ./bin

## compile: compile the app
compile: clean
	GOARCH=amd64 \
	GOOS=linux \
	CGO_ENABLED=0 \
	go build -o bin/receiver ./receiver

## test: run the unit tests
test: clean
	go test ./receiver -cover

## run: run the app locally
run:
	go run ./receiver

## build: build the docker image
build: clean
	docker build -t $(IMAGE_TAG) .

## tag: build the docker image and tag it with the repository tag
tag: build
	@echo "Tagging image with tag ${PUSH_TAG}"
	docker tag ${IMAGE_TAG} ${PUSH_TAG}

## push: build and tag the docker image and push it to the repository
push: tag
	@echo "Pushing image with tag ${PUSH_TAG}"
	docker push ${PUSH_TAG}