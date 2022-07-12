package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/dojinkimm/cryptoexchange-go"
)

func main() {
	accessKey := "<--Add Your Access Key Here-->"
	secretKey := "<--Add Your Secret Key Here-->"

	cryptoClient := crypto_exchange.NewClient(
		nil,
		crypto_exchange.WithAccessKey(accessKey),
		crypto_exchange.WithSecretKey(secretKey),
	)

	order, err := cryptoClient.UpbitService.CreateOrder("KRW-ETH", crypto_exchange.Buy, 0.001, 10000, crypto_exchange.MarketPriceBuy)
	if err != nil {
		logrus.Error(err)
	}

	fmt.Println(order)
}
