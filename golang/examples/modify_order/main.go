// Place a resting order (GTC, 3% below mid), then modify its price to 4% below.
package main

import (
	"fmt"
	"os"

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
	restPx := fmt.Sprintf("%d", int(mid*0.97))

	fmt.Printf("%s mid: $%.2f\n", coin, mid)
	fmt.Printf("Placing resting BUY %s @ %s (GTC, 3%% below mid)\n\n", sz, restPx)

	res := client.Exchange(map[string]interface{}{
		"action": map[string]interface{}{
			"type": "order",
			"orders": []interface{}{
				map[string]interface{}{
					"asset": coin,
					"side":  "buy",
					"price": restPx,
					"size":  sz,
					"tif":   "gtc",
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

	exchangeResp := result["exchangeResponse"].(map[string]interface{})
	response := exchangeResp["response"].(map[string]interface{})
	data := response["data"].(map[string]interface{})
	statuses := data["statuses"].([]interface{})

	var oid float64
	found := false
	for _, s := range statuses {
		sm := s.(map[string]interface{})
		if resting, ok := sm["resting"]; ok {
			rm := resting.(map[string]interface{})
			oid = rm["oid"].(float64)
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Could not extract OID from resting order")
		fmt.Println(client.PrettyJSON(exchangeResp))
		os.Exit(1)
	}

	newPx := fmt.Sprintf("%d", int(mid*0.96))
	fmt.Printf("Order resting (OID: %v)\n", oid)
	fmt.Printf("Modifying price: %s -> %s\n\n", restPx, newPx)

	modifyAction := map[string]interface{}{
		"type": "batchModify",
		"modifies": []interface{}{
			map[string]interface{}{
				"oid": oid,
				"order": map[string]interface{}{
					"asset": coin,
					"side":  "buy",
					"price": newPx,
					"size":  sz,
					"tif":   "gtc",
				},
			},
		},
	}

	res = client.Exchange(map[string]interface{}{"action": modifyAction})
	sig = client.SignHash(res["hash"].(string))

	modifyResult := client.Exchange(map[string]interface{}{
		"action":    modifyAction,
		"nonce":     res["nonce"],
		"signature": sig,
	})

	fmt.Println(client.PrettyJSON(modifyResult["exchangeResponse"]))
	fmt.Println("\nOrder modified.")
}
