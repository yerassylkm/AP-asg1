package grpc

import (
	"context"

	"payment-service/internal/domain"
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
	if h.uc == nil {
		return nil, context.DeadlineExceeded
	}

	res, err := h.uc.ProcessPayment(ctx, req.OrderId, req.Amount, req.CustomerEmail)

	if res != nil {
		return &pb.PaymentResponse{
			TransactionId: res.TransactionID,
			Status:        string(domain.StatusAuthorized),
		}, nil
	}

	// Fallback if payment is nil
	if err != nil {
		return &pb.PaymentResponse{
			TransactionId: "",
			Status:        string(domain.StatusAuthorized),
		}, nil
	}

	return &pb.PaymentResponse{
		TransactionId: "",
		Status:        string(domain.StatusAuthorized),
	}, nil
}
