# Deploying & Running
## Running Locally
First clone the repo:
```
git clone https://github.com/twoshark/ethereumBalanceProxy.git
```

From the root of the repo run the `build` make target.
```
make build
```
Then you can run the binary:
```
./eth_balance_proxy [COMMAND(S)] [FLAGS]
```

## Docker
### Local
```
git clone https://github.com/twoshark/ethereumBalanceProxy.git
cd ethereumBalanceProxy
docker build -t $tag
```
### Docker Hub
```
docker pull twosharks/balanceproxy:$version
```
## Kubernetes
A helm chart with a deployment, pod autoscaler and prometheus service monitor is provided at `helm/ethBalanceProxy`
### Deploying with `helm`
```
cd helm/ethBalanceProxy
helm install . -f overrideValues.yaml
```
### Deploying with `kubectl apply`
```
cd helm/ethBalanceProxy
helm template . -f overrideValues.yaml >> manifest.yaml  
kubectl apply -f manifest.yaml
```