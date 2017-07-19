Kraken Golang API Client
========================

Another Kraken golang API Client

Example usage:

```go
package main

import (
	"fmt"

	"github.com/mycroft/kraken-api"
)

func main() {
	key := ""
	secret := ""

	api := krakenapi.New(key, secret)

	pairs := [...]string{"DASHEUR", "XXBTZEUR", "XLTCZEUR", "XETCZEUR", "XETHZEUR", "XREPZEUR", "XXRPZEUR", "XZECZEUR", "XXLMZEUR", "GNOEUR", "XXMRZEUR"}
	tickers, err := api.ApiTicker(pairs[:])
	if err != nil {
		panic(err)
	}

	for pair, ticker := range tickers {
		fmt.Printf("%s: %v\n", pair, ticker)
	}

	balances, err := api.ApiBalance()
	if err != nil {
		panic(err)
	}

	for currency, balance := range balances {
		fmt.Printf("%s: %f\n", currency, balance)
	}
}
```

Notes
-----

A lot of ideas came from [kraken-go-api-client](https://github.com/beldur/kraken-go-api-client/)
