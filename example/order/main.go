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

	order, err := cryptoClient.UpbitService.CreateOrder("KRW-SOL", crypto_exchange.Buy, 1.000000, 40000.000000, crypto_exchange.Limit)
	if err != nil {
		logrus.Error(err)
	}

	fmt.Println(order)
}
