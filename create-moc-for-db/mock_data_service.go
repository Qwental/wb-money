package main

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Event представляет событие пользователя
type Event struct {
	Timestamp  time.Time
	UserID     uint64
	EventName  string
	Parameters string
}

// randomInt генерирует случайное целое число в диапазоне [min, max]
func randomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// randomPlatform возвращает случайную платформу (ios или android)
func randomPlatform() string {
	platforms := []string{"ios", "android"}
	return platforms[randomInt(0, 1)]
}

// randomPaymentMethod возвращает случайный способ оплаты
func randomPaymentMethod() string {
	methods := []string{"wallet", "card", "cash"}
	return methods[randomInt(0, 2)]
}

// generateUserEvents генерирует события для одного пользователя
func generateUserEvents(userID uint64, startDate time.Time) []Event {
	events := []Event{}
	currentDate := startDate

	// Событие: open_app
	openAppParams := fmt.Sprintf(`{"platform": "%s", "region": "RU"}`, randomPlatform())
	events = append(events, Event{
		Timestamp:  currentDate,
		UserID:     userID,
		EventName:  "open_app",
		Parameters: openAppParams,
	})

	// Случайный интервал
	currentDate = currentDate.Add(time.Second * time.Duration(randomInt(1, 60)))

	// С вероятностью 80% генерируем событие "cart"
	if rand.Float64() < 0.8 {
		cartParams := fmt.Sprintf(`{"total_amount": %d, "currency": "RUB", "n_goods": %d, "goods_list": "..."}`, randomInt(100, 10000), randomInt(1, 10))
		events = append(events, Event{
			Timestamp:  currentDate,
			UserID:     userID,
			EventName:  "cart",
			Parameters: cartParams,
		})

		// Интервал
		currentDate = currentDate.Add(time.Second * time.Duration(randomInt(1, 60)))

		// С вероятностью 70% генерируем событие "payment_methods"
		if rand.Float64() < 0.7 {
			paymentMethodsParams := fmt.Sprintf(`{"default_method": "%s"}`, randomPaymentMethod())
			events = append(events, Event{
				Timestamp:  currentDate,
				UserID:     userID,
				EventName:  "payment_methods",
				Parameters: paymentMethodsParams,
			})

			// Интервал
			currentDate = currentDate.Add(time.Second * time.Duration(randomInt(1, 60)))

			// С вероятностью 60% генерируем событие "buy"
			if rand.Float64() < 0.6 {
				buyParams := fmt.Sprintf(`{"amount": %d, "currency": "RUB", "n_goods": %d, "payment_method": "%s"}`, randomInt(100, 10000), randomInt(1, 10), randomPaymentMethod())
				events = append(events, Event{
					Timestamp:  currentDate,
					UserID:     userID,
					EventName:  "buy",
					Parameters: buyParams,
				})
			}
		}
	}

	return events
}

// generateMockData генерирует моковые данные для нескольких пользователей
func generateMockData(numUsers int, startDate time.Time) []Event {
	allEvents := []Event{}
	userID := uint64(1000) // Начальный ID пользователя

	for i := 0; i < numUsers; i++ {
		userEvents := generateUserEvents(userID, startDate)
		allEvents = append(allEvents, userEvents...)
		userID++
		startDate = startDate.Add(time.Minute * time.Duration(randomInt(1, 10)))
	}

	return allEvents
}

// insertIntoClickHouse вставляет события в ClickHouse
func insertIntoClickHouse(events []Event) error {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to ClickHouse: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()
	batch, err := conn.PrepareBatch(ctx, "INSERT INTO product_events (timestamp, user_id, event_name, parameters)")
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %v", err)
	}

	for _, event := range events {
		err := batch.Append(
			event.Timestamp,
			event.UserID,
			event.EventName,
			event.Parameters,
		)
		if err != nil {
			return fmt.Errorf("failed to append event: %v", err)
		}
	}

	return batch.Send()
}

// generateMockDataHandler обрабатывает HTTP-запрос для генерации данных
func generateMockDataHandler(w http.ResponseWriter, r *http.Request) {
	numUsersStr := r.URL.Query().Get("numUsers")
	startDateStr := r.URL.Query().Get("startDate")

	numUsers, err := strconv.Atoi(numUsersStr)
	if err != nil || numUsers <= 0 {
		numUsers = 100 // Увеличено до 100 пользователей по умолчанию
	}

	startDate, err := time.Parse("2006-01-02T15:04:05", startDateStr)
	if err != nil {
		startDate = time.Now()
	}

	events := generateMockData(numUsers, startDate)

	// Вставка в ClickHouse
	err = insertIntoClickHouse(events)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert data into ClickHouse: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Successfully inserted %d events into ClickHouse", len(events))))
}

// curl "http://localhost:3001/generate-mock-data?numUsers=50&startDate=2025-05-01T00:00:00"
func main() {
	rand.Seed(time.Now().UnixNano())

	// Получаем порт из переменной окружения, по умолчанию 3001
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	http.HandleFunc("/generate-mock-data", generateMockDataHandler)
	fmt.Printf("Mock data service running on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}
