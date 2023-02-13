# Ethereum Balance Proxy API
---

<details>
 <summary><code>GET</code> <code><b>/</b></code> </summary>

##### Parameters

> | name      |  type     | kind                  | description                                                           |
> |-----------|-----------------------|-------------------------|-----------------------------------------------------------------------|
> | None      |  required | object (JSON or YAML) | N/A  |

##### Example cURL

> ```bash
>  curl $URL/
> ```

</details>
<details>
 <summary><code>GET</code> <code><b>/ethereum/balance/:address</b></code> </summary>

##### Parameters

> | name     |  type     | kind | description       |
> |----------|------|-----------|------------------------|
> | address |  required | Path | Ethereum Wallet Address |

##### Example cURL

> ```bash
>  curl $URL/ethereum/balance/0x74630370197b4c4795bFEeF6645ee14F8cf8997D
> ```

</details>
<details>
 <summary><code>GET</code> <code><b>/ethereum/balance/:address/block/:block</b></code> </summary>

##### Parameters

> | name    |  type     | kind | description                 |
> |---------|-------------|-----------------------------|-----------------------------|
> | address |  required | Path | Ethereum Wallet Address     |
> | block   |  required | Path | Ethereum Block Number or `latest` |

##### Example cURL

> ```bash
>  curl $URL/ethereum/balance/0x74630370197b4c4795bFEeF6645ee14F8cf8997D/block/16363048
> ```

</details>