// Close a position -- worker queries your position and builds the counter-order.
package main

import (
	"fmt"

	"hyperliquid-api-examples/client"
)

const coin = "HYPE"

func main() {
	fmt.Printf("Closing %s position for %s\n\n", coin, client.Address)

	// Optional: custom slippage (default 3%, range 0.1%-10%)
	// res := client.Exchange(map[string]interface{}{"action": ..., "slippage": 0.05})  // 5% slippage

	res := client.Exchange(map[string]interface{}{
		"action": map[string]interface{}{
			"type":  "closePosition",
			"asset": coin,
			"user":  client.Address,
		},
	})

	if ctx, ok := res["closePositionContext"].(map[string]interface{}); ok {
		fmt.Printf("Position: %v %v\n", ctx["positionSize"], ctx["positionSide"])
		fmt.Printf("Close: %v %v @ %v\n", ctx["closeSide"], ctx["closeSize"], ctx["slippedPrice"])
	}

	sig := client.SignHash(res["hash"].(string))

	result := client.Exchange(map[string]interface{}{
		"action":    res["action"],
		"nonce":     res["nonce"],
		"signature": sig,
	})

	fmt.Println(client.PrettyJSON(result["exchangeResponse"]))
	fmt.Println("\nPosition closed.")
}
