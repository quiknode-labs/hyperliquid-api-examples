// Revoke builder fee approval.
package main

import (
	"fmt"

	"hyperliquid-api-examples/client"
)

func main() {
	res := client.Exchange(map[string]interface{}{
		"action": map[string]interface{}{"type": "approveBuilderFee", "maxFeeRate": "0%"},
	})
	sig := client.SignHash(res["hash"].(string))

	client.Exchange(map[string]interface{}{
		"action":    map[string]interface{}{"type": "approveBuilderFee", "maxFeeRate": "0%"},
		"nonce":     res["nonce"],
		"signature": sig,
	})

	fmt.Println("Builder fee revoked.")
}
