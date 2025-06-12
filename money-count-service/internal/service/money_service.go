//package service
//
//import (
//	"context"
//	"database/sql"
//	"encoding/json"
//	_ "encoding/json"
//	"fmt"
//	"github.com/jmoiron/sqlx"
//	"log"
//
//	pb "github.com/Qwental/wb-money/pkg/proto"
//)

package service

import (
	"context"
	"encoding/json"
	"github.com/Qwental/wb-money/pkg/proto"
	"github.com/jmoiron/sqlx"
)

type MoneyService struct {
	db *sqlx.DB
}

type BuyEvent struct {
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
}

func NewMoneyService(db *sqlx.DB) *MoneyService {
	return &MoneyService{db: db}
}

func (s *MoneyService) GetSavings(ctx context.Context, userID uint64) (*proto.GetSavingsResponse, error) {
	query := `
        SELECT parameters
        FROM product_events
        WHERE user_id = ? AND event_name = 'buy'
    `
	rows, err := s.db.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totalSavings := 0.0
	for rows.Next() {
		var params string
		if err := rows.Scan(&params); err != nil {
			return nil, err
		}

		var event BuyEvent
		if err := json.Unmarshal([]byte(params), &event); err != nil {
			return nil, err
		}

		// Если метод оплаты не кошелёк, добавляем 3% от суммы к экономии
		if event.PaymentMethod != "wallet" {
			totalSavings += event.Amount * 0.03
		}
	}

	return &proto.GetSavingsResponse{
		TotalSavings: totalSavings,
		Currency:     "RUB",
	}, nil
}
