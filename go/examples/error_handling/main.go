package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	twofinance "github.com/2Finance-Labs/2finance-sdk-client"
)

func main() {
	client := twofinance.NewFromEnv()
	if _, err := client.Analytics.ListIndicators(context.Background()); err != nil {
		var httpError *twofinance.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("request failed with status %d: %s\n", httpError.StatusCode, httpError.Body)
			return
		}
		log.Fatal(err)
	}
}
