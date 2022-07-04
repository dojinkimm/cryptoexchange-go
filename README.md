# Cryptocurrency Exchange APIs in Go
TBD

# Installation

```bash
go get -u github.com/dojinkimm/cryptoexchange-go
```

# Usage

### Public Information
```go
cryptoClient, err := crypto_exchange.NewClient(nil, crypto_exchange.Upbit)
if err != nil {
    logrus.Error(err)
}

// list current prices for given market codes
currPrices, err := cryptoClient.ListCurrentPriceByMarketCodes([]string{"KRW-BTC", "KRW-ETH", "BTC-ETH"})
if err != nil {
    logrus.Error(err)
}
```

### Private Information
```go
accessKey := "<--Add Your Access Key Here-->"
secretKey := "<--Add Your Secret Key Here-->"

cryptoClient, err := crypto_exchange.NewClient(
    nil,
    crypto_exchange.Upbit,
    crypto_exchange.WithAccessKey(accessKey),
    crypto_exchange.WithSecretKey(secretKey),
)
if err != nil {
    logrus.Error(err)
}

// list account information for given accessKey and secretKey user
accounts, err = cryptoClient.ListAccounts()
if err != nil {
    logrus.Error(err)
}
```


# Supported Cryptocurrency Exchanges

| CryptoCurrency Exchange       | REST Supported    | Websocket Support |
|-------------------------------|------------------ | ----------------- |
| [Upbit](https://upbit.com/)   | Yes               | No                |
