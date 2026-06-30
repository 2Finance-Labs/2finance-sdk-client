package main

import (
	"context"
	"fmt"
	"log"

	twofinance "github.com/2Finance-Labs/2finance-sdk-client"
	"github.com/2Finance-Labs/2finance-sdk-client/planner"
)

func main() {
	ctx := context.Background()
	client := twofinance.NewFromEnv()

	indicators, err := client.Analytics.ListIndicators(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("analytics indicators: %s\n", indicators)

	plan, err := client.Planner.TradingPlan(ctx, planner.Request{
		Goal:         "prepare a BTC rebalancing plan",
		UseAnalytics: true,
		UseTrading:   true,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("planner source: %s\n", plan.Source)
}
