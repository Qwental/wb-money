package main_test

import (
	"context"
	"github.com/Qwental/wb-money/internal/database"
	"testing"
)

func TestClickHouseQuery(t *testing.T) {
	db, err := database.NewClickHouseDB("clickhouse://localhost:9000/default")
	if err != nil {
		t.Fatalf("Failed to connect to ClickHouse: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Failed to close ClickHouseDB: %v", err)
		}
	}()

	query := `
		SELECT 
			user_id,
			JSONExtractFloat(parameters, 'amount'),
			JSONExtractString(parameters, 'payment_method')
		FROM product_events 
		WHERE event_name = 'buy'
		LIMIT 5
	`
	rows, err := db.QueryContext(context.Background(), query)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			t.Errorf("Failed to close rows: %v", err)
		}
	}()

	var count int
	for rows.Next() {
		var userID uint64
		var amount float64
		var paymentMethod string

		if err := rows.Scan(&userID, &amount, &paymentMethod); err != nil {
			t.Errorf("Scan failed: %v", err)
			continue
		}

		t.Logf("UserID: %d, Amount: %.2f, PaymentMethod: %s", userID, amount, paymentMethod)
		count++
	}

	if count == 0 {
		t.Log("Warning: query returned 0 rows")
	}

	if err := rows.Err(); err != nil {
		t.Errorf("Rows error: %v", err)
	}
}
