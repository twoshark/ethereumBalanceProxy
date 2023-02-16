# Ethereum Balance Proxy

This application proxies the `eth_getBalance` json rpc method and re-abstracts it as a REST API.

## Getting Started

```bash
make build && ./ethBalanceProxy server --upstreams "endpoint1, endpoint2"
```
```bash
go run -ldflags="-X main.Version=0.0.0 -X main.CommitHash=BAR -X main.BuildTimeStamp=BAZ" main.go server --upstreams "endpoint1, endpoint2" 
```

Further information may be found in our `docs`:
- [Rest API Methods](docs/api.md)
- [Deploying and Running the Proxy](docs/deploy.md)
- [Operational Runbook](docs/runbook.md)