## Build and Run
``` sh
cd ./producer
go build -o producer
./producer -config=config.conf

cd ./consumer
go build -o consumer
./consumer -config=config.conf
```

## Deploy to Dockerhub
When you push to master branch the github action will automatically build image and push it to my dockerhub

## Deploy to k8s
``` sh
kubectl create namespace try
cd helm
helm install --namespace=try -f values.yaml simple-kafka .
helm delete simple-kafka --namespace=try
```

---

## Install Kafka in K8S
* https://strimzi.io/quickstarts/
* You can connect to this kafka by "my-cluster-kafka-brokers.kafka:9092"
``` sh
# install
kubectl create namespace kafka
kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka
kubectl apply -f ./kafka/kafka-ephemeral-single.yaml -n kafka 
```
``` sh
# uninstall
kubectl delete -f ./kafka/kafka-ephemeral-single.yaml -n kafka
kubectl delete -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka
kubectl delete namespace kafka
```