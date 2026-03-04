// List all available markets (perps, spot, HIP-3).
package main

import (
	"fmt"

	"hyperliquid-api-examples/client"
)

func printGroup(label string, markets []interface{}) {
	if len(markets) == 0 {
		return
	}
	fmt.Printf("--- %s ---\n", label)
	limit := 10
	if len(markets) < limit {
		limit = len(markets)
	}
	for _, m := range markets[:limit] {
		market := m.(map[string]interface{})
		name, _ := market["name"].(string)
		fmt.Printf("  %-16s  index=%v  szDecimals=%v\n", name, market["index"], market["szDecimals"])
	}
	if len(markets) > 10 {
		fmt.Printf("  ... and %d more\n", len(markets)-10)
	}
	fmt.Println()
}

func main() {
	data := client.GetMarkets()

	var perps, spot []interface{}
	if p, ok := data["perps"].([]interface{}); ok {
		perps = p
	}
	if s, ok := data["spot"].([]interface{}); ok {
		spot = s
	}

	var hip3All []map[string]interface{}
	if hip3Map, ok := data["hip3"].(map[string]interface{}); ok {
		for dex, markets := range hip3Map {
			if marketList, ok := markets.([]interface{}); ok {
				for _, m := range marketList {
					market := m.(map[string]interface{})
					entry := make(map[string]interface{})
					for k, v := range market {
						entry[k] = v
					}
					entry["dex"] = dex
					hip3All = append(hip3All, entry)
				}
			}
		}
	}

	fmt.Printf("Perps: %d  |  Spot: %d  |  HIP-3: %d\n\n", len(perps), len(spot), len(hip3All))

	printGroup("Perps", perps)
	printGroup("Spot", spot)

	if len(hip3All) > 0 {
		fmt.Println("--- HIP-3 ---")
		limit := 10
		if len(hip3All) < limit {
			limit = len(hip3All)
		}
		for _, m := range hip3All[:limit] {
			display := fmt.Sprintf("%s:%s", m["dex"], m["name"])
			fmt.Printf("  %-16s  index=%v  szDecimals=%v\n", display, m["index"], m["szDecimals"])
		}
		if len(hip3All) > 10 {
			fmt.Printf("  ... and %d more\n", len(hip3All)-10)
		}
		fmt.Println()
	}
}
