# PUBLISH

## Build

> docker build -t poc-trino-api:latest .

## Sendo to Registry

```bash
docker tag poc-trino-api:latest registry-docker-registry-server.default.svc.cluster.local:5000/poc-trino-api:latest
docker push registry-docker-registry-server.default.svc.cluster.local:5000/poc-trino-api:latest
```

## Publish

```bash
kubectl apply -f infra/app-values.yaml
```

## Delete

```bash
kubectl delete -f infra/app-values.yaml
```

## Create Configuration

```bash
kubectl create configmap poc-trino-api-config -n trino --from-env-file=.env.k8s
```