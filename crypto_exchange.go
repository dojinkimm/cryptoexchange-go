package crypto_exchange

import (
	"net/http"
	"sync"
)

type CryptoExchange int

type DefaultCryptoExchangeClient struct {
	mu     sync.Mutex
	client *http.Client

	accessKey string
	secretKey string

	UpbitService *UpbitService
}

type service struct {
	client *http.Client

	baseURL string

	accessKey string
	secretKey string
}

type Option func(*DefaultCryptoExchangeClient)

// NewClient returns a new http client and cryptocurrency exchange services.
// In order to use APIs that need authentication, AccessKey and SecretKey must be provided.
func NewClient(httpClient *http.Client, opts ...Option) *DefaultCryptoExchangeClient {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	client := &DefaultCryptoExchangeClient{
		client: httpClient,
	}
	for _, opt := range opts {
		opt(client)
	}

	client.UpbitService = newUpbitService(client.copyHTTPClient(), client.accessKey, client.secretKey)

	return client
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
