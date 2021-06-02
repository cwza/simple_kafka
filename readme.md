## Build and Run
``` sh
cd ./producer
go build -o producer
ADDRESS=localhost:9092 TOPIC="my-topic" INTERVAL="1000" ./producer

cd ./consumer
go build -o consumer
ADDRESS=localhost:9092 TOPIC="my-topic" ./consumer
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