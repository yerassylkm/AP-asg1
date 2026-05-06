package grpc

import (
	"context"

	"payment-service/internal/usecase"

	pb "github.com/yerassylkm/AP-asg2_generated"
)

type PaymentGRPCHandler struct {
	pb.UnimplementedPaymentServiceServer
	uc *usecase.PaymentUseCase
}

func NewPaymentGRPCHandler(uc *usecase.PaymentUseCase) *PaymentGRPCHandler {
	return &PaymentGRPCHandler{uc: uc}
}

func (h *PaymentGRPCHandler) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	res, err := h.uc.ProcessPayment(ctx, req.OrderId, req.Amount, "")

	status := "SUCCESS"
	if err != nil {
		status = "DECLINED"
	}

	return &pb.PaymentResponse{
		TransactionId: res.TransactionID,
		Status:        status,
	}, nil
}
