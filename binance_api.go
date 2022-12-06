package binance_p2p_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type BinanceP2PApi struct {
}

func NewBinanceP2PApi() *BinanceP2PApi {
	return &BinanceP2PApi{}
}

func (b *BinanceP2PApi) GetExchange(assets string, fiat string, payTypes []string,
	tradeType string, transAmount float64) (exchange float64, minAmount float64, maxAmount float64,
	err error) {
	rawExchange, err := b.GetExchangesRaw(assets, fiat, 1, payTypes, 1, tradeType, transAmount)
	if err != nil {
		return 0, 0, 0, err
	}

	if len(rawExchange.Data) == 0 {
		return 0, 0, 0, errors.New("no results")
	}

	exchange, err = strconv.ParseFloat(rawExchange.Data[0].Adv.Price, 32)
	if err != nil {
		return 0, 0, 0, err
	}
	minAmount, err = strconv.ParseFloat(rawExchange.Data[0].Adv.MinSingleTransAmount, 32)
	if err != nil {
		return 0, 0, 0, err
	}
	maxAmount, err = strconv.ParseFloat(rawExchange.Data[0].Adv.MaxSingleTransAmount, 32)
	if err != nil {
		return 0, 0, 0, err
	}

	return exchange, minAmount, maxAmount, nil
}

func (b *BinanceP2PApi) GetExchangesRaw(assets string, fiat string, page int, payTypes []string,
	rows int, tradeType string, transAmount float64) (Response, error) {

	body := Request{
		Asset:       assets,
		Fiat:        fiat,
		Page:        page,
		PayTypes:    payTypes,
		Rows:        rows,
		TradeType:   tradeType,
		TransAmount: transAmount,
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return Response{}, err
	}

	bodyReader := bytes.NewReader(bodyJson)
	request, err := http.NewRequest("POST", bapi+getExchange, bodyReader)
	if err != nil {
		return Response{}, err
	}
	request.Header.Set(HeaderContentType, ApplicationJsonContentType)
	request.Header.Set(HeaderOrigin, P2PBinanceOrigin)
	request.Header.Set(HeaderPragma, NoCashPragma)
	request.Header.Set(HeaderTE, TrailersTE)
	request.Header.Set(HeaderUserAgent, MozillaUserAgent)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	responseRaw, err := client.Do(request)
	if err != nil {
		return Response{}, err
	}
	defer responseRaw.Body.Close()

	var response Response
	err = json.NewDecoder(responseRaw.Body).Decode(&response)
	if err != nil {
		return Response{}, err
	}
	return response, nil
}
