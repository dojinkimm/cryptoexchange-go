package service

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

const (
	upbitBaseURL = "https://api.upbit.com"
)

type UpbitService struct {
	client *http.Client

	baseURL string

	accessKey          string
	authorizationToken string
}

func NewUpbitService(httpClient *http.Client, accessKey, secretKey string) (*UpbitService, error) {
	claimMap := jwt.MapClaims{}
	claimMap["access_key"] = accessKey

	nonce, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	claimMap["nonce"] = nonce

	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, claimMap)
	token, err := claim.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &UpbitService{
		client:             httpClient,
		baseURL:            upbitBaseURL,
		accessKey:          accessKey,
		authorizationToken: "Bearer " + token,
	}, nil
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
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", s.authorizationToken)

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

	var accounts []UpbitAccount
	if err := json.Unmarshal(b, &accounts); err != nil {
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

	var marketCodes []MarketCode
	if err := json.Unmarshal(b, &marketCodes); err != nil {
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

	var currPrices []MarketCurrentPrice
	if err := json.Unmarshal(b, &currPrices); err != nil {
		return nil, err
	}

	return currPrices, nil
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
