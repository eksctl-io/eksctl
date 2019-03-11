#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

kubectl apply -f https://cloud.weave.works/k8s/scope

kubectl -n kube-system create sa tiller
kubectl create clusterrolebinding tiller-cluster-rule \
    --clusterrole=cluster-admin \
    --serviceaccount=kube-system:tiller 
helm init --upgrade --service-account=tiller --wait
kubectl apply -f https://raw.githubusercontent.com/openfaas/faas-netes/master/namespaces.yml

OPENFAAS_PASSWORD="$(head -c 12 /dev/urandom | shasum | cut -d' ' -f1)"

kubectl create secret generic basic-auth \
  --namespace=openfaas \
  --from-literal=basic-auth-user=admin \
  --from-literal=basic-a--namespace=openfaas-password=${OPENFAAS_PASSWORD}

helm upgrade openfaas --install openfaas/openfaas \
    --namespace=openfaas  \
    --set functionNamespace=openfaas-fn \
    --set serviceType=LoadBalancer \
    --set basic_auth=true \
    --set operator.create=true \
    --set gateway.replicas=4 \
    --set queueWorker.replicas=8

kubectl create namespace sock-shop
kubectl apply -f https://raw.githubusercontent.com/microservices-demo/microservices-demo/master/deploy/kubernetes/complete-demo.yaml

kubectl scale --namespace=sock-shop --replicas=20 deployment front-end
kubectl scale --namespace=sock-shop --replicas=10 deployment carts
kubectl scale --namespace=sock-shop --replicas=10 deployment orders
kubectl scale --namespace=sock-shop --replicas=10 deployment shipping
