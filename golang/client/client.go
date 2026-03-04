// Package client provides a REST client for the Hyperliquid API (via QuickNode builder API).
//
// No SDK required -- just net/http + go-ethereum.
package client

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	apiURL    = "https://send.hyperliquidapi.com"
	hlInfoURL = "https://api.hyperliquid.xyz/info"
)

var (
	privateKey *ecdsa.PrivateKey
	Address    string
	httpClient = &http.Client{}
)

func init() {
	pk := os.Getenv("PRIVATE_KEY")
	if pk == "" {
		fmt.Println("Set PRIVATE_KEY environment variable (hex, with or without 0x)")
		os.Exit(1)
	}

	pk = strings.TrimPrefix(pk, "0x")
	key, err := crypto.HexToECDSA(pk)
	if err != nil {
		fmt.Printf("Invalid PRIVATE_KEY: %v\n", err)
		os.Exit(1)
	}

	privateKey = key
	Address = crypto.PubkeyToAddress(key.PublicKey).Hex()
	fmt.Printf("Wallet: %s\n", Address)
}

// Exchange sends a POST to /exchange -- build (no signature) or send (with signature).
func Exchange(body map[string]interface{}) map[string]interface{} {
	jsonBody, _ := json.Marshal(body)
	resp, err := httpClient.Post(apiURL+"/exchange", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var data map[string]interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		fmt.Printf("Invalid JSON response: %v\n", err)
		os.Exit(1)
	}

	if errVal, ok := data["error"]; ok && errVal != nil {
		fmt.Printf("\nError (%d):\n", resp.StatusCode)
		fmt.Printf("  error:    %v\n", data["error"])
		fmt.Printf("  message:  %v\n", data["message"])
		if guidance, ok := data["guidance"]; ok && guidance != nil {
			fmt.Printf("  guidance: %v\n", guidance)
		}
		os.Exit(1)
	}

	return data
}

// GetApproval checks builder fee approval status via GET /approval?user=<addr>.
func GetApproval(user string) map[string]interface{} {
	resp, err := httpClient.Get(fmt.Sprintf("%s/approval?user=%s", apiURL, user))
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	return data
}

// GetMarkets lists all available markets via GET /markets.
func GetMarkets() map[string]interface{} {
	resp, err := httpClient.Get(apiURL + "/markets")
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	return data
}

// PostEndpoint sends a POST to a utility endpoint (e.g. /openOrders, /orderStatus).
func PostEndpoint(path string, body map[string]interface{}) map[string]interface{} {
	jsonBody, _ := json.Marshal(body)
	resp, err := httpClient.Post(apiURL+path, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	return data
}

// SignHash signs a 32-byte hash and returns {r, s, v}.
func SignHash(hashHex string) map[string]interface{} {
	hashHex = strings.TrimPrefix(hashHex, "0x")
	hashBytes, err := hex.DecodeString(hashHex)
	if err != nil {
		fmt.Printf("Invalid hash hex: %v\n", err)
		os.Exit(1)
	}

	sig, err := crypto.Sign(hashBytes, privateKey)
	if err != nil {
		fmt.Printf("Signing failed: %v\n", err)
		os.Exit(1)
	}

	// sig is 65 bytes: R (32) || S (32) || V (1)
	r := hex.EncodeToString(sig[:32])
	s := hex.EncodeToString(sig[32:64])
	v := int(sig[64]) + 27

	return map[string]interface{}{
		"r": "0x" + r,
		"s": "0x" + s,
		"v": v,
	}
}

// GetMid gets the current mid price for a coin from Hyperliquid.
func GetMid(coin string) float64 {
	body, _ := json.Marshal(map[string]string{"type": "allMids"})
	resp, err := httpClient.Post(hlInfoURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)

	if val, ok := data[coin]; ok {
		if s, ok := val.(string); ok {
			f, err := strconv.ParseFloat(s, 64)
			if err == nil {
				return f
			}
		}
	}
	return 0
}

// GetHip3Mid gets the mid price for a HIP-3 market (requires dex parameter).
func GetHip3Mid(coin string) float64 {
	dex := strings.Split(coin, ":")[0]
	body, _ := json.Marshal(map[string]interface{}{"type": "allMids", "dex": dex})
	resp, err := httpClient.Post(hlInfoURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)

	if val, ok := data[coin]; ok {
		if s, ok := val.(string); ok {
			f, err := strconv.ParseFloat(s, 64)
			if err == nil {
				return f
			}
		}
	}
	return 0
}

// PrettyJSON marshals a value to indented JSON string.
func PrettyJSON(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}
