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

	accounts, err := cryptoClient.UpbitService.ListAccounts()
	if err != nil {
		logrus.Error(err)
	}

	for _, acc := range accounts {
		fmt.Println(acc)
	}
}
