#!/usr/bin/env bash

podname=$(kubectl get po -o name -n balance-proxy-test --context kind-balance-proxy-cluster)
kubectl port-forward --context kind-balance-proxy-cluster "${podname}" 8080:8080 &
echo "Wait for forwarding.."
sleep 20
sh scripts/debug/checkApi.sh