# Operational Runbook

## Grafana

In the event of an issue or an alert, the grafana dashboard provides a snapshot of important metrics to help guide diagnosis
and support healthy operations.

## RPC Calls 
```
ethBalanceProxy ops call eth_syncing|eth_getBlockNumber --endpoint $ETH_RPC_URL
ethBalanceProxy ops call eth_getBalance --endpoint $ETH_RPC_URL --address $wallet_address --block $block
```
Bypassing the proxy server and upstream manager, this will call utilizes the `upstream.ethereum.Client` to execute
ethereum json rpc calls in the same fashion as the upstream manager

## Debug Scripts
The shell scripts in `scripts/debug/` are included to support debugging and development. They are not intended for
production use.

### checkApi.sh
Performs a quick call to each endpoint produced by the api

## Ops commands
These commands are intended to be used to aid in debugging, ci and ops scripts. They are not intended for production use.

### Health Check
```
ethBalanceProxy ops healthCheck --endpoint $ETH_RPC_URL
```
Directly calls `upstream.ethereum.Client{}.HealthCheck` against a single endpoint

## Unit Tests
```make test```

## `kind` testing

To support local testing in `kind`, a set of `kind:...` make targets are provided.

### `make kind:init`

Instantiates kind cluster and installs prometheus operator & crds

### `make kind:install` / `make kind:uninstall`

Installs/uninstalls the helm chart to/from a kind cluster created with `make kind:init`

### `make kind:test`

Forwards `kind` pod port 8080 to local 8080 and runs `scripts/debug/checkApi.sh` to exercise all the endpoints in the proxy

### `make kind:teardown`

Teardown `kind` cluster