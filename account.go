package crypto_exchange

import (
	"errors"
	"strconv"
)

const (
	currencyKRW = "KRW"
)

type Account struct {
	Currency              string  `json:"currency,omitempty"`
	BalanceQuantity       float64 `json:"balance_quantity,omitempty"`
	LockedQuantity        float64 `json:"locked_quantity,omitempty"`
	AveragePurchaseAmount float64 `json:"average_purchase_amount,omitempty"`
	UnitCurrency          string  `json:"unit_currency,omitempty"`
}

var (
	ErrUnsupportedCryptoExchange = errors.New("unsupported crypto exchange")
)

// ListAccounts fetches user's account information in cryptocurrency exchange
func (c *DefaultCryptoExchangeClient) ListAccounts() ([]*Account, error) {
	var accounts []*Account
	if c.cryptoExchange == Upbit {
		upbitAccounts, err := c.upbitService.ListAccounts()
		if err != nil {
			return nil, err
		}

		for _, acc := range upbitAccounts {
			// currency KRW is actually a deposit
			if acc.Currency == currencyKRW {
				continue
			}
			balance, err := strconv.ParseFloat(acc.Balance, 64)
			if err != nil {
				return nil, err
			}

			locked, err := strconv.ParseFloat(acc.Locked, 64)
			if err != nil {
				return nil, err
			}

			avgBuyPrice, err := strconv.ParseFloat(acc.AvgBuyPrice, 64)
			if err != nil {
				return nil, err
			}

			accounts = append(accounts, &Account{
				Currency:              acc.Currency,
				BalanceQuantity:       balance,
				LockedQuantity:        locked,
				AveragePurchaseAmount: avgBuyPrice,
				UnitCurrency:          acc.UnitCurrency,
			})
		}

		return accounts, nil
	}

	return nil, ErrUnsupportedCryptoExchange
}
