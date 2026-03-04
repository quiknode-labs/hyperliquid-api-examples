// Place a market order -- no price needed, worker computes it automatically.
package main

import (
	"fmt"

	"hyperliquid-api-examples/client"
)

const (
	coin = "BTC"
	size = "0.00011"
)

func main() {
	fmt.Printf("Market BUY %s %s\n\n", size, coin)

	// Optional: custom slippage (default 3%, range 0.1%-10%)
	// res := client.Exchange(map[string]interface{}{"action": ..., "slippage": 0.05})  // 5% slippage

	res := client.Exchange(map[string]interface{}{
		"action": map[string]interface{}{
			"type": "order",
			"orders": []interface{}{
				map[string]interface{}{
					"asset": coin,
					"side":  "buy",
					"size":  size,
					"tif":   "market",
				},
			},
		},
	})

	action := res["action"].(map[string]interface{})
	orders := action["orders"].([]interface{})
	firstOrder := orders[0].(map[string]interface{})
	fmt.Printf("Computed price (mid + slippage, default 3%%): %v\n", firstOrder["p"])
	fmt.Printf("Builder fee: %v\n", res["builderFee"])

	sig := client.SignHash(res["hash"].(string))

	result := client.Exchange(map[string]interface{}{
		"action":    res["action"],
		"nonce":     res["nonce"],
		"signature": sig,
	})

	fmt.Println(client.PrettyJSON(result["exchangeResponse"]))
}
