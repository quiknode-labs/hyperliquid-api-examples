// Check what happened to an order by OID.
package main

import (
	"fmt"
	"os"
	"strconv"

	"hyperliquid-api-examples/client"
)

func main() {
	oid := 0
	if len(os.Args) > 1 {
		oid, _ = strconv.Atoi(os.Args[1])
	}
	if oid == 0 {
		fmt.Println("Usage: go run . <oid>")
		fmt.Println("Get OIDs from: go run ../open_orders")
		os.Exit(1)
	}

	result := client.PostEndpoint("/orderStatus", map[string]interface{}{
		"user": client.Address,
		"oid":  oid,
	})

	status, _ := result["status"].(string)
	explanation, _ := result["explanation"].(string)

	if status == "unknownOid" {
		fmt.Printf("Order %d: not found\n", oid)
		fmt.Printf("  %s\n", explanation)
		os.Exit(0)
	}

	name, _ := result["name"].(string)
	spot := ""
	if isSpot, ok := result["isSpot"].(bool); ok && isSpot {
		spot = " [SPOT]"
	}

	fmt.Printf("Order %d on %s%s: %s\n", oid, name, spot, status)
	fmt.Printf("  %s\n", explanation)
	if order, ok := result["order"]; ok && order != nil {
		fmt.Printf("  Details: %s\n", client.PrettyJSON(order))
	}
}
