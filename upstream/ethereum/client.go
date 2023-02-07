package ethereum

import "github.com/ethereum/go-ethereum/ethclient"

type IClient interface {
	Connect() error
}

type Client struct {
	endpoint  string
	ethClient *ethclient.Client
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

func (c *Client) Connect() error {
	ethClient, err := ethclient.Dial(c.endpoint)
	if err != nil {
		return err
	}
	c.ethClient = ethClient
}
