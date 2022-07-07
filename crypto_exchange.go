package crypto_exchange

import (
	"net/http"
	"sync"

	"github.com/dojinkimm/cryptoexchange-go/service"
)

//go:generate mockgen -package=client -destination=crypto_exchange_mock.go . CryptoExchangeClient
type CryptoExchangeClient interface {
	ListAccounts() ([]*Account, error)
	ListTradableMarkets() ([]*Market, error)
	ListCurrentPriceByMarketCodes([]string) ([]*MarketCurrentPrice, error)
}

type CryptoExchange int

// Supported Crypto Currency Exchanges
const (
	// Upbit format
	Upbit CryptoExchange = iota + 1
)

type DefaultCryptoExchangeClient struct {
	mu     sync.Mutex
	client *http.Client

	accessKey      string
	secretKey      string
	cryptoExchange CryptoExchange

	upbitService *service.UpbitService
}

type Option func(*DefaultCryptoExchangeClient)

// NewClient returns a new http client and cryptocurrency exchange services.
// In order to use APIs that need authentication, AccessKey and SecretKey must be provided.
func NewClient(httpClient *http.Client, cryptoExchange CryptoExchange, opts ...Option) (*DefaultCryptoExchangeClient, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	client := &DefaultCryptoExchangeClient{
		client:         httpClient,
		cryptoExchange: cryptoExchange,
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
	return func(c *DefaultCryptoExchangeClient) {
		c.accessKey = accessKey
	}
}

func WithSecretKey(secretKey string) Option {
	return func(c *DefaultCryptoExchangeClient) {
		c.secretKey = secretKey
	}
}

// copyHTTPClient returns the client use by DefaultClient
func (c *DefaultCryptoExchangeClient) copyHTTPClient() *http.Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	copiedClient := *c.client
	return &copiedClient
}
