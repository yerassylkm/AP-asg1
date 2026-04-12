package grpc

import (
	"context" // <--- ДОБАВЬ ЭТУ СТРОКУ

	pb "github.com/yerassylkm/AP-asg2_generated"
	"payment-service/internal/usecase"
)

type PaymentGRPCHandler struct {
    pb.UnimplementedPaymentServiceServer
    // Добавь звездочку *, так как PaymentUseCase — это структура
    uc *usecase.PaymentUseCase 
}

// Здесь тоже добавь звездочку перед usecase.PaymentUseCase
func NewPaymentGRPCHandler(uc *usecase.PaymentUseCase) *PaymentGRPCHandler {
    return &PaymentGRPCHandler{uc: uc}
}

// Внутри самого метода ProcessPayment вызови свой реальный метод:
func (h *PaymentGRPCHandler) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
    // Вызываем твой метод из payment_uc.go
    res, err := h.uc.ProcessPayment(ctx, req.OrderId, req.Amount)
    
    status := "SUCCESS"
    if err != nil {
        status = "DECLINED"
    }

    return &pb.PaymentResponse{
        TransactionId: res.TransactionID,
        Status:        status,
        // CreatedAt заполни по желанию через timestamppb.Now()
    }, nil
}