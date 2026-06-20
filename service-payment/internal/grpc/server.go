// ─── gRPC Server Implementation ───────────────────────────────
package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	nexusflow "github.com/MikeCHOKKI/nexusflow/protos/gen/go/nexusflow"

	"github.com/MikeCHOKKI/nexusflow/service-payment/internal/model"
	"github.com/MikeCHOKKI/nexusflow/service-payment/internal/payment"
)

// PaymentServer implements the protobuf PaymentServiceServer interface.
type PaymentServer struct {
	nexusflow.UnimplementedPaymentServiceServer
	provider *payment.StubProvider
}

// NewPaymentServer creates a PaymentServer backed by the given stub provider.
func NewPaymentServer(provider *payment.StubProvider) *PaymentServer {
	return &PaymentServer{provider: provider}
}

// ProcessPayment handles an incoming payment request.
func (s *PaymentServer) ProcessPayment(
	ctx context.Context,
	req *nexusflow.ProcessPaymentRequest,
) (*nexusflow.PaymentResponse, error) {
	// ── Validate input ──────────────────────────────────────────
	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "amount must be positive")
	}
	if req.Currency == "" {
		return nil, status.Error(codes.InvalidArgument, "currency is required")
	}
	if req.Method == nexusflow.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "payment method is required")
	}

	// ── Process ─────────────────────────────────────────────────
	payment, err := s.provider.ProcessPayment(
		ctx,
		req.OrderId,
		req.UserId,
		req.Amount,
		req.Currency,
		model.PaymentMethod(req.Method),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "process payment: %v", err)
	}

	// ── Response ────────────────────────────────────────────────
	msg := "Payment processed successfully"
	if payment.Status == model.Failed {
		msg = "Payment failed"
	}

	return &nexusflow.PaymentResponse{
		Payment: payment.ToProto(),
		Message: msg,
	}, nil
}

// GetPayment retrieves a payment by ID.
func (s *PaymentServer) GetPayment(
	ctx context.Context,
	req *nexusflow.GetPaymentRequest,
) (*nexusflow.Payment, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "payment id is required")
	}

	p, err := s.provider.GetPayment(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "payment not found: %v", err)
	}

	return p.ToProto(), nil
}

// ListPayments returns a paginated list of payments.
func (s *PaymentServer) ListPayments(
	ctx context.Context,
	req *nexusflow.ListPaymentsRequest,
) (*nexusflow.ListPaymentsResponse, error) {
	page := int(req.GetPage())
	if page < 1 {
		page = 1
	}
	limit := int(req.GetLimit())
	if limit < 1 || limit > 100 {
		limit = 20
	}

	statusFilter := model.PaymentStatus(-1)
	if req.GetStatus() != nexusflow.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED {
		statusFilter = model.PaymentStatus(req.GetStatus())
	}

	payments, total, err := s.provider.ListPayments(ctx, page, limit, statusFilter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list payments: %v", err)
	}

	protos := make([]*nexusflow.Payment, 0, len(payments))
	for _, p := range payments {
		protos = append(protos, p.ToProto())
	}

	totalPages := (total + limit - 1) / limit

	return &nexusflow.ListPaymentsResponse{
		Payments: protos,
		Pagination: &nexusflow.Pagination{
			Total:      int32(total),
			Page:       int32(page),
			Limit:      int32(limit),
			TotalPages: int32(totalPages),
		},
	}, nil
}

// RefundPayment processes a refund for an existing payment.
func (s *PaymentServer) RefundPayment(
	ctx context.Context,
	req *nexusflow.RefundPaymentRequest,
) (*nexusflow.PaymentResponse, error) {
	if req.PaymentId == "" {
		return nil, status.Error(codes.InvalidArgument, "payment_id is required")
	}

	p, err := s.provider.RefundPayment(ctx, req.PaymentId, req.Reason)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "refund failed: %v", err)
	}

	return &nexusflow.PaymentResponse{
		Payment: p.ToProto(),
		Message: "Payment refunded successfully",
	}, nil
}
