# indexer
Indexer knative service

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
