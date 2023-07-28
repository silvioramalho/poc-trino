.PHONY: docker-build docker-push deploy

IMAGE_NAME := poc-trino-api
IMAGE_TAG := latest
REGISTRY := registry-docker-registry-server.default.svc.cluster.local:5000

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

docker-push: docker-build
	docker push $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

deploy:
	kubectl apply -f infra/app-values.yaml

all: docker-build docker-push deploy
