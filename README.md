# Cryptocurrency Exchange APIs in Go
TBD

# Installation

```bash
go get -u github.com/dojinkimm/cryptoexchange-go
```

# Usage

### Public Information
```go
cryptoClient := crypto_exchange.NewClient(nil)

// list current prices for given market codes
currPrices, err := cryptoClient.UpbitService.ListCurrentPriceByMarketCodes([]string{"KRW-BTC", "KRW-ETH", "BTC-ETH"})
if err != nil {
    logrus.Error(err)
}
```

### Private Information
```go
accessKey := "<--Add Your Access Key Here-->"
secretKey := "<--Add Your Secret Key Here-->"

cryptoClient := crypto_exchange.NewClient(
    nil,
    crypto_exchange.WithAccessKey(accessKey),
    crypto_exchange.WithSecretKey(secretKey),
)

// list account information for given accessKey and secretKey user
accounts, err := cryptoClient.UpbitService.ListAccounts()
if err != nil {
    logrus.Error(err)
}
```


# Supported Cryptocurrency Exchanges

| Cryptocurrency Exchange     | REST Supported    | Websocket Support |
|-----------------------------|------------------ | ----------------- |
| [Upbit](https://upbit.com/) | Yes               | No                |

# Supported APIs

| APIs                             | Upbit |
|----------------------------------|-------|
| List All Accounts                | ✅     |
| List Tradable Market Codes       | ✅     |
| Get Current Price of Market Code | ✅     |
| Order                            | ✅     |
