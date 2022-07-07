package crypto_exchange

import (
	"time"
)

type Change int

const (
	Even Change = iota + 1
	Rise
	Fall
)

type Market struct {
	// ex) KRW-BTC, BTC-ETC
	MarketCode  string `json:"market,omitempty"`
	KoreanName  string `json:"korean_name,omitempty"`
	EnglishName string `json:"english_name,omitempty"`
	IsRisky     bool   `json:"is_risky,omitempty"`
}

type MarketCurrentPrice struct {
	// ex) KRW-BTC, BTC-ETC
	MarketCode                string    `json:"market_code,omitempty"`
	TradedAtUTC               time.Time `json:"traded_at_utc,omitempty"`
	TradedAtMilliseconds      int64     `json:"traded_at_milliseconds,omitempty"`
	OpeningPrice              float64   `json:"opening_price,omitempty"`
	HighPrice                 float64   `json:"high_price,omitempty"`
	LowPrice                  float64   `json:"low_price,omitempty"`
	TradePrice                float64   `json:"trade_price,omitempty"`
	PreviousClosingPrice      float64   `json:"previous_closing_price,omitempty"`
	Change                    Change    `json:"change,omitempty"`
	ChangePrice               float64   `json:"change_price,omitempty"`
	ChangeRate                float64   `json:"change_rate,omitempty"`
	SignedChangePrice         float64   `json:"signed_change_price,omitempty"`
	SignedChangeRate          float64   `json:"signed_change_rate,omitempty"`
	TradeVolume               float64   `json:"trade_volume,omitempty"`
	AccumulatedTradePrice     float64   `json:"accumulated_trade_price,omitempty"`
	AccumulatedTradePrice24H  float64   `json:"accumulated_trade_price_24h,omitempty"`
	AccumulatedTradeVolume    float64   `json:"accumulated_trade_volume,omitempty"`
	AccumulatedTradeVolume24H float64   `json:"accumulated_trade_volume_24h,omitempty"`
	Highest52WeekPrice        float64   `json:"highest_52_week_price,omitempty"`
	Highest52WeekDate         string    `json:"highest_52_week_date,omitempty"`
	Lowest52WeekPrice         float64   `json:"lowest_52_week_price,omitempty"`
	Lowest52WeekDate          string    `json:"lowest_52_week_date,omitempty"`
	TimestampMilliseconds     int64     `json:"timestamp_milliseconds,omitempty"`
}

var (
	riskMessageMap = map[string]bool{
		"NONE":    false,
		"CAUTION": true,
	}

	changeMap = map[string]Change{
		"EVEN": Even,
		"RISE": Rise,
		"FALL": Fall,
	}
)

// ListTradableMarkets fetches markets that can be traded in cryptocurrency exchange
func (c *DefaultCryptoExchangeClient) ListTradableMarkets() ([]*Market, error) {
	var marketCodes []*Market
	if c.cryptoExchange == Upbit {
		upbitMarketCodes, err := c.upbitService.ListMarketCodes()
		if err != nil {
			return nil, err
		}

		for _, mc := range upbitMarketCodes {
			marketCodes = append(marketCodes, &Market{
				MarketCode:  mc.Market,
				KoreanName:  mc.KoreanName,
				EnglishName: mc.EnglishName,
				IsRisky:     riskMessageMap[mc.MarketWarning],
			})
		}

		return marketCodes, nil
	}

	return nil, ErrUnsupportedCryptoExchange
}

// ListCurrentPriceByMarketCodes fetches current price of a market by market codes
func (c *DefaultCryptoExchangeClient) ListCurrentPriceByMarketCodes(marketCodes []string) ([]*MarketCurrentPrice, error) {
	var marketCurrPrices []*MarketCurrentPrice
	if c.cryptoExchange == Upbit {
		upbitCurrPrices, err := c.upbitService.ListCurrentPriceByMarketCodes(marketCodes)
		if err != nil {
			return nil, err
		}

		for _, p := range upbitCurrPrices {
			marketCurrPrices = append(marketCurrPrices, &MarketCurrentPrice{
				MarketCode:                p.Market,
				TradedAtUTC:               time.UnixMilli(p.TradeTimestamp).UTC(),
				TradedAtMilliseconds:      p.TradeTimestamp,
				OpeningPrice:              p.OpeningPrice,
				HighPrice:                 p.HighPrice,
				LowPrice:                  p.LowPrice,
				TradePrice:                p.TradePrice,
				PreviousClosingPrice:      p.PrevClosingPrice,
				Change:                    changeMap[p.Change],
				ChangePrice:               p.ChangePrice,
				ChangeRate:                p.ChangeRate,
				SignedChangePrice:         p.SignedChangePrice,
				SignedChangeRate:          p.SignedChangeRate,
				TradeVolume:               p.TradeVolume,
				AccumulatedTradePrice:     p.AccTradePrice,
				AccumulatedTradePrice24H:  p.AccTradePrice24H,
				AccumulatedTradeVolume:    p.AccTradeVolume,
				AccumulatedTradeVolume24H: p.AccTradeVolume24H,
				Highest52WeekPrice:        p.Highest52WeekPrice,
				Highest52WeekDate:         p.Highest52WeekDate,
				Lowest52WeekPrice:         p.Lowest52WeekPrice,
				Lowest52WeekDate:          p.Lowest52WeekDate,
				TimestampMilliseconds:     p.Timestamp,
			})
		}

		return marketCurrPrices, nil
	}

	return nil, ErrUnsupportedCryptoExchange
}
