package main

import (
	"context"
	"log"
	"os"

	twofinance "github.com/2Finance-Labs/2finance-sdk-client"
	"github.com/2Finance-Labs/2finance-sdk-client/auth"
)

func main() {
	ctx := context.Background()
	tokenSource, err := auth.NewClientCredentialsTokenSource(auth.ClientCredentialsConfig{
		TokenURL:     os.Getenv("TWO_FINANCE_AUTH_TOKEN_URL"),
		ClientID:     os.Getenv("TWO_FINANCE_AUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWO_FINANCE_AUTH_CLIENT_SECRET"),
		Scopes:       []string{"2finance.sdk"},
	})
	if err != nil {
		log.Fatal(err)
	}

	client := twofinance.New(twofinance.Config{
		AnalyticsURL: os.Getenv("TWO_FINANCE_ANALYTICS_URL"),
		TokenSource:  tokenSource,
	})
	if _, err := client.Analytics.ListIndicators(ctx); err != nil {
		log.Fatal(err)
	}
}
