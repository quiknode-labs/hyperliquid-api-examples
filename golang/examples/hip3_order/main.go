// Place an IOC order on a HIP-3 market by name (e.g. xyz:SILVER).
package main

import (
	"fmt"
	"strconv"

	"hyperliquid-api-examples/client"
)

const coin = "xyz:SILVER"

func main() {
	mid := client.GetHip3Mid(coin)
	if mid == 0 {
		fmt.Printf("Could not fetch %s mid price, using fallback\n", coin)
		mid = 78.0
	}

	sz := fmt.Sprintf("%.2f", 11.0/mid)
	buyPx := fmt.Sprintf("%.2f", mid*1.03)

	szFloat, _ := strconv.ParseFloat(sz, 64)
	fmt.Printf("%s mid: $%.2f\n", coin, mid)
	fmt.Printf("BUY %s @ %s (IOC, ~$%.2f notional)\n\n", sz, buyPx, szFloat*mid)

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
