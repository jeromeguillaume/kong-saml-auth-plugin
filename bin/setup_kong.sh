#!/bin/sh

source env.sh

# kubectl create namespace $NAMESPACE

echo "apiVersion: v1
kind: Namespace
metadata:
  name: $NAMESPACE" | kubectl apply -f -

#create ConfigMap for JS plugins
pushd ../kong-js-pdk
kubectl create configmap kong-plugin-js --from-file=plugins -n $NAMESPACE
popd

# create the Kong EE license secret
kubectl create secret generic kong-enterprise-license --from-file=./license -n $NAMESPACE

echo "Deploying Kong EE..."
pushd ../deploy
helm install $RELEASE kong/kong \
--values values.yaml -n $NAMESPACE \
--set ingressController.installCRDs=false

popd
echo "Waiting for Kong EE to start..."