package binance_p2p_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type BinanceP2PApi struct {
	client *http.Client
}

func NewBinanceP2PApi(client *http.Client) *BinanceP2PApi {
	return &BinanceP2PApi{
		client: client,
	}
}

func (b *BinanceP2PApi) GetExchange(assets string, fiat string, payTypes []string,
	tradeType TradeType, transAmount float64) (Exchange, error) {
	rawExchange, err := b.GetExchangesRaw(assets, fiat, 1, payTypes, 1, tradeType, transAmount)
	if err != nil {
		return Exchange{}, err
	}

	if len(rawExchange.Data) == 0 {
		return Exchange{}, errors.New("no results")
	}

	price, err := strconv.ParseFloat(rawExchange.Data[0].Adv.Price, 32)
	if err != nil {
		return Exchange{}, err
	}
	minAmount, err := strconv.ParseFloat(rawExchange.Data[0].Adv.MinSingleTransAmount, 32)
	if err != nil {
		return Exchange{}, err
	}
	maxAmount, err := strconv.ParseFloat(rawExchange.Data[0].Adv.MaxSingleTransAmount, 32)
	if err != nil {
		return Exchange{}, err
	}

	return Exchange{
		Price:     price,
		MinAmount: minAmount,
		MaxAmount: maxAmount,
	}, nil
}

func (b *BinanceP2PApi) GetExchangesRaw(assets string, fiat string, page int, payTypes []string,
	rows int, tradeType TradeType, transAmount float64) (Response, error) {

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

	responseRaw, err := b.client.Do(request)
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
