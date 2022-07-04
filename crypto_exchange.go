package crypto_exchange

import (
	"net/http"
	"sync"

	"github.com/dojinkimm/cryptoexchange-go/service"
)

type Client interface {
	ListAccounts() ([]*Account, error)
	ListTradableMarkets() ([]*Market, error)
	ListCurrentPriceByMarketCodes() ([]string, error)
}

type CryptoCurrencyExchange int

// Supported Crypto Currency Exchanges
const (
	// Upbit format
	Upbit CryptoCurrencyExchange = iota + 1
)

type DefaultClient struct {
	mu     sync.Mutex
	client *http.Client

	accessKey              string
	secretKey              string
	cryptocurrencyExchange CryptoCurrencyExchange

	upbitService *service.UpbitService
}

type Option func(*DefaultClient)

// NewClient returns a new http client and cryptocurrency exchange services.
// In order to use APIs that need authentication, AccessKey and SecretKey must be provided.
func NewClient(httpClient *http.Client, cryptocurrencyExchange CryptoCurrencyExchange, opts ...Option) (*DefaultClient, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	client := &DefaultClient{
		client:                 httpClient,
		cryptocurrencyExchange: cryptocurrencyExchange,
	}
	for _, opt := range opts {
		opt(client)
	}

	var err error
	client.upbitService, err = service.NewUpbitService(client.copyHTTPClient(), client.accessKey, client.secretKey)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func WithAccessKey(accessKey string) Option {
	return func(c *DefaultClient) {
		c.accessKey = accessKey
	}
}

func WithSecretKey(secretKey string) Option {
	return func(c *DefaultClient) {
		c.secretKey = secretKey
	}
}

// copyHTTPClient returns the client use by DefaultClient
func (c *DefaultClient) copyHTTPClient() *http.Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	copiedClient := *c.client
	return &copiedClient
}
