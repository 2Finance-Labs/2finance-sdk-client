package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	twofinance "github.com/2Finance-Labs/2finance-sdk-client"
)

func main() {
	ctx := context.Background()
	client := twofinance.NewFromEnv()

	var response json.RawMessage
	err := client.Analytics.PostWithOptions(
		ctx,
		"/analytics/candles:upsert",
		map[string]any{"symbol": "BTC-USDT"},
		&response,
		twofinance.WithHeader("X-Trace-ID", "trace-1"),
		twofinance.WithIdempotencyKey("candles-upsert-001"),
		twofinance.WithQueryParam("source", "sdk-example"),
		twofinance.WithPagination(1, 25),
		twofinance.WithTimeout(5*time.Second),
		twofinance.WithMaxRetries(1),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("response: %s\n", response)
}
