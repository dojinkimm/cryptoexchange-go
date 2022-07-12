package crypto_exchange

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type UpbitService service

const (
	upbitBaseURL = "https://api.upbit.com"
)

func newUpbitService(httpClient *http.Client, accessKey, secretKey string) *UpbitService {
	return &UpbitService{
		client:    httpClient,
		baseURL:   upbitBaseURL,
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

type UpbitAccount struct {
	Currency            string `json:"currency,omitempty"`
	Balance             string `json:"balance,omitempty"`
	Locked              string `json:"locked,omitempty"`
	AvgBuyPrice         string `json:"avg_buy_price,omitempty"`
	AvgBuyPriceModified bool   `json:"avg_buy_price_modified,omitempty"`
	UnitCurrency        string `json:"unit_currency,omitempty"`
}

type UpbitErrorBody struct {
	Message string `json:"message,omitempty"`
	Name    string `json:"name,omitempty"`
}

type UpbitError struct {
	Error UpbitErrorBody `json:"error,omitempty"`
}

func (s *UpbitService) ListAccounts() ([]UpbitAccount, error) {
	req, err := http.NewRequest(http.MethodGet, s.baseURL+"/v1/accounts", nil)
	if err != nil {
		return nil, err
	}

	authToken, err := generateAuthorizationToken(s.accessKey, s.secretKey, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", authToken)

	resp, err := getResponse(req)
	if err != nil {
		return nil, err
	}

	var accounts []UpbitAccount
	if err := json.Unmarshal(resp, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

type MarketCode struct {
	// ex) KRW-BTC, BTC-ETC
	Market      string `json:"market,omitempty"`
	KoreanName  string `json:"korean_name,omitempty"`
	EnglishName string `json:"english_name,omitempty"`
	// NONE (해당 사항 없음), CAUTION(투자유의)
	MarketWarning string `json:"market_warning,omitempty"`
}

func (s *UpbitService) ListMarketCodes() ([]MarketCode, error) {
	req, err := http.NewRequest(http.MethodGet, s.baseURL+"/v1/market/all", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("isDetails", "true")
	req.URL.RawQuery = q.Encode()

	resp, err := getResponse(req)
	if err != nil {
		return nil, err
	}

	var marketCodes []MarketCode
	if err := json.Unmarshal(resp, &marketCodes); err != nil {
		return nil, err
	}

	return marketCodes, nil
}

type MarketCurrentPrice struct {
	// ex) KRW-BTC, BTC-ETC
	Market             string  `json:"market,omitempty"`
	TradeDate          string  `json:"trade_date,omitempty"`
	TradeTime          string  `json:"trade_time,omitempty"`
	TradeDateKst       string  `json:"trade_date_kst,omitempty"`
	TradeTimeKst       string  `json:"trade_time_kst,omitempty"`
	TradeTimestamp     int64   `json:"trade_timestamp,omitempty"`
	OpeningPrice       float64 `json:"opening_price,omitempty"`
	HighPrice          float64 `json:"high_price,omitempty"`
	LowPrice           float64 `json:"low_price,omitempty"`
	TradePrice         float64 `json:"trade_price,omitempty"`
	PrevClosingPrice   float64 `json:"prev_closing_price,omitempty"`
	Change             string  `json:"change,omitempty"`
	ChangePrice        float64 `json:"change_price,omitempty"`
	ChangeRate         float64 `json:"change_rate,omitempty"`
	SignedChangePrice  float64 `json:"signed_change_price,omitempty"`
	SignedChangeRate   float64 `json:"signed_change_rate,omitempty"`
	TradeVolume        float64 `json:"trade_volume,omitempty"`
	AccTradePrice      float64 `json:"acc_trade_price,omitempty"`
	AccTradePrice24H   float64 `json:"acc_trade_price_24h,omitempty"`
	AccTradeVolume     float64 `json:"acc_trade_volume,omitempty"`
	AccTradeVolume24H  float64 `json:"acc_trade_volume_24h,omitempty"`
	Highest52WeekPrice float64 `json:"highest_52_week_price,omitempty"`
	Highest52WeekDate  string  `json:"highest_52_week_date,omitempty"`
	Lowest52WeekPrice  float64 `json:"lowest_52_week_price,omitempty"`
	Lowest52WeekDate   string  `json:"lowest_52_week_date,omitempty"`
	Timestamp          int64   `json:"timestamp,omitempty"`
}

func (s *UpbitService) ListCurrentPriceByMarketCodes(marketCodes []string) ([]MarketCurrentPrice, error) {
	req, err := http.NewRequest(http.MethodGet, s.baseURL+"/v1/ticker", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("markets", strings.Join(marketCodes, ","))
	req.URL.RawQuery = q.Encode()

	resp, err := getResponse(req)
	if err != nil {
		return nil, err
	}

	var currPrices []MarketCurrentPrice
	if err := json.Unmarshal(resp, currPrices); err != nil {
		return nil, err
	}

	return currPrices, nil
}

type BuyOrSell string
type OrderType string

const (
	Buy  BuyOrSell = "bid"
	Sell BuyOrSell = "ask"

	Limit           OrderType = "limit"
	MarketPriceBuy  OrderType = "price"
	MarketPriceSell OrderType = "market"
)

type OrderResp struct {
	Uuid            string    `json:"uuid,omitempty"`
	Side            string    `json:"side,omitempty"`
	OrdType         string    `json:"ord_type,omitempty"`
	Price           string    `json:"price,omitempty"`
	AvgPrice        string    `json:"avg_price,omitempty"`
	State           string    `json:"state,omitempty"`
	Market          string    `json:"market,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	Volume          string    `json:"volume,omitempty"`
	RemainingVolume string    `json:"remaining_volume,omitempty"`
	ReservedFee     string    `json:"reserved_fee,omitempty"`
	RemainingFee    string    `json:"remaining_fee,omitempty"`
	PaidFee         string    `json:"paid_fee,omitempty"`
	Locked          string    `json:"locked,omitempty"`
	ExecutedVolume  string    `json:"executed_volume,omitempty"`
	TradesCount     int       `json:"trades_count,omitempty"`
}

func (s *UpbitService) CreateOrder(
	marketCode string,
	side BuyOrSell,
	volume float64,
	price float64,
	orderType OrderType,
) (*OrderResp, error) {
	params := url.Values{}
	params.Add("market", marketCode)
	params.Add("side", string(side))
	params.Add("volume", fmt.Sprintf("%f", volume))
	params.Add("price", fmt.Sprintf("%f", price))
	params.Add("ord_type", string(orderType))
	encodedParams := params.Encode()

	jsonPayload := []byte(fmt.Sprintf(`{"market": "%s","side": "%s","volume": "%s","price": "%s","ord_type": "%s"}`,
		marketCode, string(side), fmt.Sprintf("%f", volume), fmt.Sprintf("%f", price), string(orderType)))
	req, err := http.NewRequest(http.MethodPost, s.baseURL+"/v1/orders", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	authHeader, err := generateAuthorizationToken(s.accessKey, s.secretKey, &encodedParams)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", authHeader)

	resp, err := getResponse(req)
	if err != nil {
		return nil, err
	}

	var orderResp *OrderResp
	if err := json.Unmarshal(resp, orderResp); err != nil {
		return nil, err
	}

	return orderResp, nil
}

func unmarshalUpbitError(b []byte, code int) error {
	var upbitError UpbitError
	if err := json.Unmarshal(b, &upbitError); err != nil {
		return err
	}

	return errors.New("status code " + strconv.Itoa(code) + ": " + upbitError.Error.Message)
}

func closeBody(resp *http.Response) func() {
	return func() {
		if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
			logrus.Error(err.Error())
		}

		if err := resp.Body.Close(); err != nil {
			logrus.Error(err.Error())
		}
	}
}

func getResponse(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer closeBody(resp)
	}
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, unmarshalUpbitError(b, resp.StatusCode)
	}

	return b, nil
}
