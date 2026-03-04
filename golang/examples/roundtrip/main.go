// IOC buy then sell on BTC perp -- full trade cycle (~$11 notional).
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"hyperliquid-api-examples/client"
)

const coin = "BTC"

func sendOrder(side, px, sz string) map[string]interface{} {
	res := client.Exchange(map[string]interface{}{
		"action": map[string]interface{}{
			"type": "order",
			"orders": []interface{}{
				map[string]interface{}{
					"asset": coin,
					"side":  side,
					"price": px,
					"size":  sz,
					"tif":   "ioc",
				},
			},
		},
	})
	sig := client.SignHash(res["hash"].(string))

	return client.Exchange(map[string]interface{}{
		"action":    res["action"],
		"nonce":     res["nonce"],
		"signature": sig,
	})
}

func checkStatuses(resp map[string]interface{}, label string) bool {
	response, ok := resp["response"].(map[string]interface{})
	if !ok {
		return true
	}
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return true
	}
	statuses, ok := data["statuses"].([]interface{})
	if !ok {
		return true
	}
	for _, s := range statuses {
		if sm, ok := s.(map[string]interface{}); ok {
			if errVal, ok := sm["error"]; ok {
				fmt.Printf("%s error: %v\n", label, errVal)
				return false
			}
		}
	}
	return true
}

func main() {
	mid := client.GetMid(coin)
	if mid == 0 {
		fmt.Printf("Could not fetch %s mid price\n", coin)
		os.Exit(1)
	}

	sz := fmt.Sprintf("%.5f", 11.0/mid)
	szFloat, _ := strconv.ParseFloat(sz, 64)
	fmt.Printf("%s mid: $%.2f\n", coin, mid)
	fmt.Printf("Trade size: %s %s (~$%.2f)\n\n", sz, coin, szFloat*mid)

	buyPx := fmt.Sprintf("%d", int(mid*1.03))
	fmt.Printf("BUY %s @ %s (IOC)\n", sz, buyPx)
	buyResult := sendOrder("buy", buyPx, sz)
	buyResp := buyResult["exchangeResponse"].(map[string]interface{})
	if !checkStatuses(buyResp, "Buy") {
		os.Exit(1)
	}
	fmt.Printf("Buy filled: %s\n\n", client.PrettyJSON(buyResp))

	time.Sleep(1 * time.Second)

	sellPx := fmt.Sprintf("%d", int(mid*0.97))
	fmt.Printf("SELL %s @ %s (IOC)\n", sz, sellPx)
	sellResult := sendOrder("sell", sellPx, sz)
	sellResp := sellResult["exchangeResponse"].(map[string]interface{})
	if !checkStatuses(sellResp, "Sell") {
		os.Exit(1)
	}
	fmt.Printf("Sell filled: %s\n", client.PrettyJSON(sellResp))

	fmt.Println("\nRound-trip complete. Position should be flat.")
}
