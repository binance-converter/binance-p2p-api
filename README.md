# binance-p2p-api
## Usage
```go
client := http.Client{Timeout: 5 * time.Second}
binance := binance_p2p_api.NewBinanceP2PApi(&client)

var totalRUB float64 = 5000
exchange, err := binance.GetExchange("USDT", "RUB", []string{"QIWI"}, binance_p2p_api.TradeTypeBuy, totalRUB)
if err != nil {
    log.Fatal(err.Error())
}
fmt.Printf("%+v\n", exchange)
```