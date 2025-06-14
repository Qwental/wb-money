package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"net/http"
	"os"
)

type AnalyticsResult struct {
	WalletPaymentMethodsShare float64 `json:"wallet_payment_methods_share"`
	WalletPurchaseShare       float64 `json:"wallet_purchase_share"`
}

func getAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to ClickHouse: %v", err), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	ctx := context.Background()

	query1 := `
		SELECT
			COUNT(DISTINCT CASE WHEN JSONExtractString(parameters, 'default_method') = 'wallet' THEN user_id END) /
			COUNT(DISTINCT user_id) AS wallet_payment_methods_share
		FROM product_events
		WHERE event_name = 'payment_methods'
	`

	var walletPaymentMethodsShare float64
	err = conn.QueryRow(ctx, query1).Scan(&walletPaymentMethodsShare)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute query 1: %v", err), http.StatusInternalServerError)
		return
	}

	query2 := `
		SELECT
			COUNT(DISTINCT CASE WHEN buy.user_id IS NOT NULL THEN buy.user_id END) /
			COUNT(DISTINCT open_app.user_id) AS wallet_purchase_share
		FROM
			(SELECT DISTINCT user_id FROM product_events WHERE event_name = 'open_app') AS open_app
		LEFT JOIN
			(SELECT DISTINCT user_id
			 FROM product_events
			 WHERE event_name = 'buy'
			 AND JSONExtractString(parameters, 'payment_method') = 'wallet') AS buy
		ON open_app.user_id = buy.user_id
	`

	var walletPurchaseShare float64
	err = conn.QueryRow(ctx, query2).Scan(&walletPurchaseShare)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute query 2: %v", err), http.StatusInternalServerError)
		return
	}

	result := AnalyticsResult{
		WalletPaymentMethodsShare: walletPaymentMethodsShare,
		WalletPurchaseShare:       walletPurchaseShare,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}

	http.HandleFunc("/analytics", getAnalyticsHandler)
	fmt.Printf("Analytics service running on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}

// curl http://localhost:3002/analytics
