package main

import (
	"fmt"
	"github.com/Qwental/wb-money/internal/database"
	"github.com/Qwental/wb-money/internal/handler"
	"github.com/Qwental/wb-money/internal/service"
	"github.com/Qwental/wb-money/pkg/proto"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"os"

	"github.com/improbable-eng/grpc-web/go/grpcweb"

	"log"
	"net"
	"net/http"
)

// Значения по умолчанию
const (
	defaultGrpcPort = ":50051"
	defaultWebPort  = ":8080"
	defaultHost     = "0.0.0.0"
)

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// buildDSN создает строку подключения к ClickHouse из переменных окружения
func buildDSN() string {
	host := getEnv("CLICKHOUSE_HOST", "localhost")
	port := getEnv("CLICKHOUSE_PORT", "9000")
	db := getEnv("CLICKHOUSE_DB", "default")
	user := getEnv("CLICKHOUSE_USER", "default")
	password := getEnv("CLICKHOUSE_PASSWORD", "")

	if password != "" {
		return fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s", user, password, host, port, db)
	}
	return fmt.Sprintf("clickhouse://%s@%s:%s/%s", user, host, port, db)
}

func main() {
	// Конфигурация из переменных окружения
	grpcHost := getEnv("GRPC_HOST", defaultHost)
	grpcPort := getEnv("GRPC_PORT", defaultGrpcPort)
	webPort := getEnv("WEB_PORT", defaultWebPort)

	grpcAddr := fmt.Sprintf("%s:%s", grpcHost, grpcPort)
	webAddr := fmt.Sprintf("%s:%s", grpcHost, webPort)

	dsn := buildDSN()

	log.Printf("Starting Money Service...")
	log.Printf("DSN: %s", dsn)
	log.Printf("gRPC server will start on: %s", grpcAddr)
	log.Printf("gRPC-Web server will start on: %s", webAddr)

	// Подключение к ClickHouse
	db, err := database.NewClickHouseDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close ClickHouseDB: %v", err)
		}
	}(db)

	log.Printf("Successfully connected to ClickHouse")

	// Инициализация сервисов
	svc := service.NewMoneyService(db)
	h := handler.NewMoneyHandler(svc)

	// Создание gRPC сервера
	grpcServer := grpc.NewServer()
	proto.RegisterMoneyServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	// gRPC-Web обёртка
	wrappedGrpc := grpcweb.WrapServer(grpcServer,
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(origin string) bool {
			// В продакшене здесь должна быть более строгая проверка
			return true
		}),
	)

	// HTTP сервер для gRPC-Web
	httpServer := &http.Server{
		Addr: webAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// CORS заголовки для preflight запросов
			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Grpc-Web, X-User-Agent")
				w.WriteHeader(http.StatusOK)
				return
			}

			if wrappedGrpc.IsGrpcWebRequest(r) || wrappedGrpc.IsAcceptableGrpcCorsRequest(r) {
				wrappedGrpc.ServeHTTP(w, r)
			} else {
				// Для обычных HTTP запросов можно добавить health check
				if r.URL.Path == "/health" {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte("OK"))
					if err != nil {
						return

					}
					return
				}
				http.NotFound(w, r)
			}
		}),
	}

	// Запускаем обычный gRPC сервер в отдельной горутине
	go func() {
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("Failed to listen on %s: %v", grpcAddr, err)
		}
		log.Printf("gRPC server started on %s", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Запускаем gRPC-Web сервер
	log.Printf("gRPC-Web server started on %s", webAddr)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
