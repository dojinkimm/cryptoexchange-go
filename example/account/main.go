package main

import (
	"github.com/sirupsen/logrus"

	"github.com/dojinkimm/cryptoexchange-go"
)

func main() {
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

	_, err = cryptoClient.ListAccounts()
	if err != nil {
		logrus.Error(err)
	}
}
