TAG?=latest
REGISTRY?=gcr.io/tmogoserverless

build:
	docker build -t storefinder/indexer .

push:
	docker tag storefinder/indexer:$(TAG) $(REGISTRY)/storefinder/indexer:$(TAG)
	docker push $(REGISTRY)/storefinder/indexer:$(TAG)