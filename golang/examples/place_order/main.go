// Place a perp limit order (BTC, IOC, ~$10 notional).
package main

import (
	"fmt"
	"os"
	"strconv"

	"hyperliquid-api-examples/client"
)

const coin = "BTC"

func main() {
	mid := client.GetMid(coin)
	if mid == 0 {
		fmt.Printf("Could not fetch %s mid price\n", coin)
		os.Exit(1)
	}

	sz := fmt.Sprintf("%.5f", 11.0/mid)
	buyPx := fmt.Sprintf("%d", int(mid*1.03))

	szFloat, _ := strconv.ParseFloat(sz, 64)
	fmt.Printf("%s mid: $%.2f\n", coin, mid)
	fmt.Printf("BUY %s @ %s (IOC, ~$%.2f notional)\n", sz, buyPx, szFloat*mid)

	res := client.Exchange(map[string]interface{}{
		"action": map[string]interface{}{
			"type": "order",
			"orders": []interface{}{
				map[string]interface{}{
					"asset": coin,
					"side":  "buy",
					"price": buyPx,
					"size":  sz,
					"tif":   "ioc",
				},
			},
		},
	})
	sig := client.SignHash(res["hash"].(string))

	result := client.Exchange(map[string]interface{}{
		"action":    res["action"],
		"nonce":     res["nonce"],
		"signature": sig,
	})

	fmt.Println(client.PrettyJSON(result["exchangeResponse"]))
}
