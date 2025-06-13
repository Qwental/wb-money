package handler

import (
	"context"
	"log"

	"github.com/Qwental/wb-money/internal/service"
	"github.com/Qwental/wb-money/pkg/proto"
)

type MoneyHandler struct {
	proto.UnimplementedMoneyServiceServer
	svc *service.MoneyService
}

func NewMoneyHandler(svc *service.MoneyService) *MoneyHandler {
	return &MoneyHandler{svc: svc}
}

func (h *MoneyHandler) GetSavings(ctx context.Context, req *proto.GetSavingsRequest) (*proto.GetSavingsResponse, error) {
	// Логируем входящий запрос
	log.Printf("- запрос GetSavings для пользователя: %d", req.UserId)

	// Валидация запроса
	if req.UserId <= 0 {
		log.Printf("Некорректный User ID: %d", req.UserId)
		return &proto.GetSavingsResponse{
			Status:  proto.GetSavingsResponse_INVALID_REQUEST,
			Message: "User ID должен быть положительным числом",
		}, nil
	}

	// Вызываем сервис
	response, err := h.svc.GetSavings(ctx, uint64(req.UserId))
	if err != nil {
		log.Printf("Ошибка сервиса для пользователя %d: %v", req.UserId, err)
		return &proto.GetSavingsResponse{
			Status:  proto.GetSavingsResponse_UNKNOWN_ERROR,
			Message: "Внутренняя ошибка сервера",
		}, nil
	}

	// Логируем результат
	log.Printf("- Результат для пользователя %d: статус=%s, экономия=%.2f, покупок=%d, wb=%d",
		req.UserId, response.Status.String(), response.TotalSavings, response.TotalPurchases, response.WbCardPurchases)

	return response, nil
}
