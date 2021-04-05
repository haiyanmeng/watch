#!/bin/bash

docker build -t gcr.io/haiyanmeng-anthos/watch:latest .
docker push gcr.io/haiyanmeng-anthos/watch
kubectl apply -f config.yaml
