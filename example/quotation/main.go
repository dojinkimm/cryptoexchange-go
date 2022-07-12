package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/dojinkimm/cryptoexchange-go"
)

func main() {
	cryptoClient := crypto_exchange.NewClient(nil)
	markets, err := cryptoClient.UpbitService.ListMarketCodes()
	if err != nil {
		logrus.Error(err)
	}

	for _, m := range markets {
		fmt.Println(m)
		break
	}

	currPrices, err := cryptoClient.UpbitService.ListCurrentPriceByMarketCodes([]string{"KRW-BTC", "KRW-ETH", "BTC-ETH"})
	if err != nil {
		logrus.Error(err)
	}

	for _, m := range currPrices {
		fmt.Println(m)
		break
	}
}
