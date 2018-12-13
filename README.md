# indexer
Indexer knative service

## Build indexer service using Knative build
At this point we should be able to build the query service. Since the build for query leverages a shareable knative CI build template be sure to create it before running the build.

```
kubectl apply -f build.yaml
```

### Checking build
Run command below to check the status of build

```
kubectl get builds -o yaml
```

## Deployment

### Create Service account in GCP 

###Assign Pubsub editor role

### Download the credentials JSON

### Create secret

```
kubectl create secret generic gcppubsub-source-key --from-file=key.json={path to service account credentials json}

```

### Deploy indexer service

```
kubectl apply -f service.yaml

```

### Deploy Gcppubsource


```
kubectl apply -f gcppubsub-source.yaml

```
