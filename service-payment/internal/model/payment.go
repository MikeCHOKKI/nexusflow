// ─── Domain Models ─────────────────────────────────────────────
package model

import (
	"time"

	nexusflow "github.com/MikeCHOKKI/nexusflow/protos/gen/go/nexusflow"
)

// PaymentMethod maps to the proto PaymentMethod enum.
type PaymentMethod int32

const (
	PaymentMethodUnspecified PaymentMethod = 0
	OrangeMoney              PaymentMethod = 1
	MTNMoMo                  PaymentMethod = 2
	Wave                     PaymentMethod = 3
	Visa                     PaymentMethod = 4
	Mastercard               PaymentMethod = 5
)

// PaymentStatus maps to the proto PaymentStatus enum.
type PaymentStatus int32

const (
	PaymentStatusUnspecified PaymentStatus = 0
	Initiated                PaymentStatus = 1
	Processing               PaymentStatus = 2
	Completed                PaymentStatus = 3
	Failed                   PaymentStatus = 4
	Refunded                 PaymentStatus = 5
)

// Payment is the internal domain representation of a payment transaction.
type Payment struct {
	ID             string
	OrderID        string
	UserID         string
	Amount         float64
	Currency       string
	Method         PaymentMethod
	Status         PaymentStatus
	TransactionRef string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ToProto converts an internal Payment to the protobuf Payment message.
func (p *Payment) ToProto() *nexusflow.Payment {
	return &nexusflow.Payment{
		Id:             p.ID,
		OrderId:        p.OrderID,
		UserId:         p.UserID,
		Amount:         p.Amount,
		Currency:       p.Currency,
		Method:         nexusflow.PaymentMethod(p.Method),
		Status:         nexusflow.PaymentStatus(p.Status),
		TransactionRef: p.TransactionRef,
		Metadata: &nexusflow.Metadata{
			CreatedAt: p.CreatedAt.Format(time.RFC3339),
			UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
		},
	}
}

// FromProtoPayment converts a protobuf Payment to the internal domain model.
func FromProtoPayment(p *nexusflow.Payment) *Payment {
	return &Payment{
		ID:             p.Id,
		OrderID:        p.OrderId,
		UserID:         p.UserId,
		Amount:         p.Amount,
		Currency:       p.Currency,
		Method:         PaymentMethod(p.Method),
		Status:         PaymentStatus(p.Status),
		TransactionRef: p.TransactionRef,
	}
}

// MethodLabel returns a human-readable label for the payment method.
func (m PaymentMethod) Label() string {
	switch m {
	case OrangeMoney:
		return "Orange Money"
	case MTNMoMo:
		return "MTN MoMo"
	case Wave:
		return "Wave"
	case Visa:
		return "Visa"
	case Mastercard:
		return "Mastercard"
	default:
		return "Unknown"
	}
}

// StatusLabel returns a human-readable label for the payment status.
func (s PaymentStatus) Label() string {
	switch s {
	case Initiated:
		return "Initiated"
	case Processing:
		return "Processing"
	case Completed:
		return "Completed"
	case Failed:
		return "Failed"
	case Refunded:
		return "Refunded"
	default:
		return "Unknown"
	}
}

// ProtoMethod converts a domain PaymentMethod to the proto enum.
func (m PaymentMethod) Proto() nexusflow.PaymentMethod {
	return nexusflow.PaymentMethod(m)
}
