# Operational Runbook

## Deploying & Running
### Running Locally
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

### Docker
#### Local
```
git clone https://github.com/twoshark/ethereumBalanceProxy.git
cd ethereumBalanceProxy
docker build -t $tag
```
#### Docker Hub
```
docker pull twosharks/balanceproxy:$version
```
### Kubernetes
A helm chart with a deployment, pod autoscaler and prometheus service monitor is provided at `helm/ethBalanceProxy`
#### Deploying with `helm`
```
cd helm/ethBalanceProxy
helm install . -f overrideValues.yaml
```
#### Deploying with `kubectl apply`
```
cd helm/ethBalanceProxy
helm template . -f overrideValues.yaml >> manifest.yaml  
kubectl apply -f manifest.yaml
```

# Testing, Developing & Debugging
## Unit Tests
```make test```

## Debug Scripts
The shell scripts in `scripts/debug/` are included to support debugging and development. They are not intended for 
production use.

### checkApi.sh
Performs a quick call to each endpoint produced by the api

## Ops commands
These commands are intended to be used to aid in debugging, ci and ops scripts. They are not intended for production use.

### RPC Calls 
```
ethBalanceProxy ops call eth_syncing|eth_getBlockNumber --endpoint $ETH_RPC_URL
ethBalanceProxy ops call eth_getBalance --endpoint $ETH_RPC_URL --address $wallet_address --block $block
```
Bypassing the proxy server and upstream manager, this will call utilizes the `upstream.ethereum.Client` to execute
ethereum json rpc calls in the same fashion as the upstream manager

### Health Check
```
ethBalanceProxy ops healthCheck --endpoint $ETH_RPC_URL
```
Directly calls `upstream.ethereum.Client{}.HealthCheck` against a single endpoint


