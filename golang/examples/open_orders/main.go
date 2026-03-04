// View open orders with enriched asset info and pre-built cancel actions.
package main

import (
	"fmt"

	"hyperliquid-api-examples/client"
)

func main() {
	result := client.PostEndpoint("/openOrders", map[string]interface{}{"user": client.Address})

	count := int(result["count"].(float64))
	fmt.Printf("Open orders: %d\n", count)

	if orders, ok := result["orders"].([]interface{}); ok {
		for _, o := range orders {
			order := o.(map[string]interface{})
			side := "SELL"
			if order["side"].(string) == "B" {
				side = "BUY"
			}
			spot := ""
			if isSpot, ok := order["isSpot"].(bool); ok && isSpot {
				spot = " [SPOT]"
			}
			fmt.Printf("  %s%s %s %s @ %s (OID: %v)\n",
				order["name"], spot, side, order["sz"], order["limitPx"], order["oid"])
		}
	}

	if count > 0 {
		cancelActions := result["cancelActions"].(map[string]interface{})

		fmt.Println("\nCancel actions by asset:")
		byAsset := cancelActions["byAsset"].(map[string]interface{})
		for name, action := range byAsset {
			a := action.(map[string]interface{})
			cancels := a["cancels"].([]interface{})
			fmt.Printf("  %s: %d order(s) — pass as action to POST /exchange\n", name, len(cancels))
		}

		fmt.Println("\nTo cancel ALL orders, pass cancelActions.all as the action to POST /exchange:")
		fmt.Println(client.PrettyJSON(cancelActions["all"]))
	}
}
