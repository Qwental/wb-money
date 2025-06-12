package main

import (
	"github.com/Qwental/wb-money/internal/database"
	"github.com/Qwental/wb-money/internal/handler"
	"github.com/Qwental/wb-money/internal/service"
	"github.com/Qwental/wb-money/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/improbable-eng/grpc-web/go/grpcweb"

	"log"
	"net"
	"net/http"
)

const (
	grpcPort = ":50051" // оригинальный (tcp)
	webPort  = ":8080"  // для фронта
	dsn      = "clickhouse://localhost:9000/default"
)

func main() {
	db, err := database.NewClickHouseDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}
	defer db.Close()

	svc := service.NewMoneyService(db)
	h := handler.NewMoneyHandler(svc)

	grpcServer := grpc.NewServer()
	proto.RegisterMoneyServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	// 👇 gRPC-Web обёртка
	wrappedGrpc := grpcweb.WrapServer(grpcServer)

	// 👇 HTTP сервер
	httpServer := &http.Server{
		Addr: webPort,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if wrappedGrpc.IsGrpcWebRequest(r) || wrappedGrpc.IsAcceptableGrpcCorsRequest(r) {
				wrappedGrpc.ServeHTTP(w, r)
			} else {
				http.NotFound(w, r)
			}
		}),
	}

	// 🎧 Запускаем слушание TCP для обычного gRPC (для отладки и grpcurl)
	go func() {
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		log.Printf("gRPC server (pure) started on %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	log.Printf("gRPC-Web server started on %s", webPort)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
