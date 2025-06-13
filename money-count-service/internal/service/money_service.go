package service

import (
	"context"
	_ "database/sql"
	"encoding/json"
	"fmt"
	"log"

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
	// Валидация входных данных
	if userID == 0 {
		return &proto.GetSavingsResponse{
			Status:  proto.GetSavingsResponse_INVALID_REQUEST,
			Message: "User ID не может быть пустым",
		}, nil
	}

	// Проверяем существование пользователя
	var userExists bool
	checkUserQuery := `SELECT count(*) > 0 FROM product_events WHERE user_id = ?`
	err := s.db.GetContext(ctx, &userExists, checkUserQuery, userID)
	if err != nil {
		log.Printf("Ошибка проверки пользователя %d: %v", userID, err)
		return &proto.GetSavingsResponse{
			Status:  proto.GetSavingsResponse_DB_ERROR,
			Message: "Ошибка доступа к базе данных",
		}, nil
	}

	if !userExists {
		return &proto.GetSavingsResponse{
			Status:  proto.GetSavingsResponse_USER_NOT_FOUND,
			Message: fmt.Sprintf("Пользователь с ID %d не найден", userID),
		}, nil
	}

	// Получаем данные о покупках
	query := `
        SELECT parameters
        FROM product_events
        WHERE user_id = ? AND event_name = 'buy'
    `
	rows, err := s.db.QueryxContext(ctx, query, userID)
	if err != nil {
		log.Printf("Ошибка выполнения запроса для пользователя %d: %v", userID, err)
		return &proto.GetSavingsResponse{
			Status:  proto.GetSavingsResponse_DB_ERROR,
			Message: "Ошибка получения данных о покупках",
		}, nil
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("error closing rows: %v", err)
		}
	}(rows)

	var (
		totalSavings      float64
		totalPurchases    int32
		wbCardPurchases   int32
		hasValidPurchases bool
	)

	for rows.Next() {
		var params string
		if err := rows.Scan(&params); err != nil {
			log.Printf("Ошибка сканирования данных для пользователя %d: %v", userID, err)
			continue
		}

		var event BuyEvent
		if err := json.Unmarshal([]byte(params), &event); err != nil {
			log.Printf("Ошибка парсинга JSON для пользователя %d: %v", userID, err)
			continue
		}

		// Увеличиваем счетчик покупок
		totalPurchases++
		hasValidPurchases = true

		// Если метод оплаты не кошелёк, добавляем 3% от суммы к экономии
		if event.PaymentMethod != "wallet" {
			totalSavings += event.Amount * 0.03
			wbCardPurchases++
		}
	}

	// Проверяем ошибки при итерации
	if err := rows.Err(); err != nil {
		log.Printf("Ошибка итерации по результатам для пользователя %d: %v", userID, err)
		return &proto.GetSavingsResponse{
			Status:  proto.GetSavingsResponse_DB_ERROR,
			Message: "Ошибка обработки данных",
		}, nil
	}

	// Если у пользователя нет покупок
	if !hasValidPurchases {
		return &proto.GetSavingsResponse{
			Status:          proto.GetSavingsResponse_NO_PURCHASES,
			TotalSavings:    0,
			Currency:        "RUB",
			TotalPurchases:  0,
			WbCardPurchases: 0,
			Message:         "У пользователя нет покупок",
		}, nil
	}

	// Формируем сообщение
	var message string
	if totalSavings > 0 {
		message = fmt.Sprintf("Вы сэкономили %.2f ₽ благодаря WB Card! Покупок с картой: %d из %d",
			totalSavings, wbCardPurchases, totalPurchases)
	} else {
		message = fmt.Sprintf("Пока нет экономии. Используйте WB Card для получения 3%% кэшбека! Всего покупок: %d",
			totalPurchases)
	}

	return &proto.GetSavingsResponse{
		Status:          proto.GetSavingsResponse_OK,
		TotalSavings:    totalSavings,
		Currency:        "RUB",
		TotalPurchases:  totalPurchases,
		WbCardPurchases: wbCardPurchases,
		Message:         message,
	}, nil
}
