# poc-trino


## Kafka Cluster 

### Strimzi

```bash
helm repo add strimzi https://strimzi.io/charts/
helm install strimzi-operator strimzi/strimzi-kafka-operator -n trino --create-namespace
```

### Kafka

```bash
kubectl apply -f kafka/kafka-cluster.yaml -n trino
```

Verify installation

```
kubectl -n trino get kafka my-kafka-cluster
```

Create Topic

```bash
kubectl -n trino get kafkatopic my-topic
```

Producer for test

```bash
kubectl -n trino run kafka-producer -ti --image=quay.io/strimzi/kafka:0.34.0-kafka-3.4.0 --rm=true --restart=Never -- bin/kafka-console-producer.sh --broker-list my-kafka-cluster-kafka-bootstrap:9092 --topic my-topic
```

Consumer for test

```bash
kubectl -n trino run kafka-consumer -ti --image=quay.io/strimzi/kafka:0.34.0-kafka-3.4.0 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-kafka-cluster-kafka-bootstrap:9092 --topic my-topic --from-beginning
```

Delete Topic

```bash
kubectl -n trino delete kafkatopic my-topic
```

Deploy producer

```bash
kubectl apply -f apps/producer.yaml -n trino
```

Deploy Topics COnfig

```bash
kubectl apply -f kafka/connect-status-topic.yaml
kubectl apply -f kafka/connect-configs-topic.yaml
kubectl apply -f kafka/connect-offsets-topic.yaml
```
Deploy Kafka-Connect

```bash
kubectl apply -f kafka/connect.yaml
```
Deploy S3 Connector

```bash
kubectl apply -f kafka/connector.yaml
```



helm install benthos-webhook benthos/benthos --values apps/helm-values-httpServer-kafkaOutput.yaml -n trino






# k8s

Install Ingress

```
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

helm install ingress-nginx ingress-nginx/ingress-nginx
```