package handler

import (
	"context"
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
	return h.svc.GetSavings(ctx, uint64(req.UserId))
}
