package grpc

import (
	"context"
	"os"

	pb "github.com/yerassylkm/AP-asg2_generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentClient struct {
	client pb.PaymentServiceClient
}

func NewPaymentClient() (*PaymentClient, error) {
	addr := os.Getenv("PAYMENT_SERVICE_ADDR")
	if addr == "" {
		addr = "payment-service:50051" 
	}

	conn, err :=grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &PaymentClient{
		client: pb.NewPaymentServiceClient(conn),
	}, nil
}

func (c *PaymentClient) ProcessPayment(ctx context.Context, orderID string, amount int64) (string, error) {
	resp, err := c.client.ProcessPayment(ctx, &pb.PaymentRequest{
		OrderId: orderID,
		Amount:  amount,
	})
	if err != nil {
		return "Failed", err
	}
	return resp.Status, nil
}