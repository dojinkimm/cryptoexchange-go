package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/dojinkimm/cryptoexchange-go"
)

func main() {
	cryptoClient, err := crypto_exchange.NewClient(nil, crypto_exchange.Upbit)
	if err != nil {
		logrus.Error(err)
	}

	markets, err := cryptoClient.ListTradableMarkets()
	if err != nil {
		logrus.Error(err)
	}

	for _, m := range markets {
		fmt.Println(m)
		break
	}

	currPrices, err := cryptoClient.ListCurrentPriceByMarketCodes([]string{"KRW-BTC", "KRW-ETH", "BTC-ETH"})
	if err != nil {
		logrus.Error(err)
	}

	for _, m := range currPrices {
		fmt.Println(m)
		break
	}
}
