// Check builder fee approval status.
package main

import (
	"fmt"

	"hyperliquid-api-examples/client"
)

func main() {
	res := client.GetApproval(client.Address)
	fmt.Println(client.PrettyJSON(res))
}
