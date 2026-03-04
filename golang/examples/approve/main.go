// Approve builder fee (one-time setup).
package main

import (
	"fmt"

	"hyperliquid-api-examples/client"
)

const maxFee = "1%"

func main() {
	res := client.Exchange(map[string]interface{}{
		"action": map[string]interface{}{"type": "approveBuilderFee", "maxFeeRate": maxFee},
	})
	sig := client.SignHash(res["hash"].(string))

	client.Exchange(map[string]interface{}{
		"action":    map[string]interface{}{"type": "approveBuilderFee", "maxFeeRate": maxFee},
		"nonce":     res["nonce"],
		"signature": sig,
	})

	fmt.Println("Builder fee approved.")
}
