// ─── Payment Processor (Stub) ──────────────────────────────────
package payment

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/MikeCHOKKI/nexusflow/service-payment/internal/db"
	"github.com/MikeCHOKKI/nexusflow/service-payment/internal/model"
)

// StubProvider simulates a payment gateway (e.g. Stripe, Orange Money API).
//
// Behaviour:
//   - 90 % of payments complete successfully, 10 % fail.
//   - Transaction references follow the format "PROV-{unix_nano}-{random}".
//   - All operations are thread-safe via sync.Mutex.
type StubProvider struct {
	mu       sync.Mutex
	rng      *rand.Rand
	store    *db.RedisClient
	payments map[string]*model.Payment // in-memory fallback alongside Redis
}

// NewStubProvider creates a new StubProvider with the given Redis store.
func NewStubProvider(store *db.RedisClient) *StubProvider {
	return &StubProvider{
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
		store:    store,
		payments: make(map[string]*model.Payment),
	}
}

// ProcessPayment simulates processing a payment through an external provider.
// It returns a domain Payment with a generated ID, status, and transaction ref.
func (sp *StubProvider) ProcessPayment(
	_ context.Context,
	orderID, userID string,
	amount float64,
	currency string,
	method model.PaymentMethod,
) (*model.Payment, error) {

	sp.mu.Lock()
	defer sp.mu.Unlock()

	now := time.Now().UTC()

	// --- Simulate processing delay (50–200 ms) ---
	time.Sleep(time.Duration(50+sp.rng.Intn(150)) * time.Millisecond)

	// --- Decide outcome: 90 % success, 10 % failure ---
	var status model.PaymentStatus
	if sp.rng.Float64() < 0.9 {
		status = model.Completed
	} else {
		status = model.Failed
	}

	// --- Generate transaction reference ---
	ref := fmt.Sprintf("PROV-%d-%06d", now.UnixNano(), sp.rng.Intn(999999))

	payment := &model.Payment{
		ID:             uuid.NewString(),
		OrderID:        orderID,
		UserID:         userID,
		Amount:         amount,
		Currency:       currency,
		Method:         method,
		Status:         status,
		TransactionRef: ref,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// Persist to Redis and in-memory fallback
	if sp.store != nil {
		if err := sp.store.SavePayment(context.Background(), payment); err != nil {
			return nil, fmt.Errorf("persist payment: %w", err)
		}
	}
	sp.payments[payment.ID] = payment

	return payment, nil
}

// GetPayment retrieves a stored payment by ID. It checks Redis first,
// then falls back to the in-memory map.
func (sp *StubProvider) GetPayment(ctx context.Context, id string) (*model.Payment, error) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	// Try Redis first
	if sp.store != nil {
		p, err := sp.store.GetPayment(ctx, id)
		if err != nil {
			return nil, err
		}
		if p != nil {
			return p, nil
		}
	}

	// Fallback to in-memory
	p, ok := sp.payments[id]
	if !ok {
		return nil, fmt.Errorf("payment %s not found", id)
	}
	return p, nil
}

// ListPayments returns all stored payments (paginated).
func (sp *StubProvider) ListPayments(ctx context.Context, page, limit int, statusFilter model.PaymentStatus) ([]*model.Payment, int, error) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	// Collect all payments from in-memory store
	var all []*model.Payment
	for _, p := range sp.payments {
		all = append(all, p)
	}

	// Sort by created_at descending
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})

	// Apply status filter if specified
	if statusFilter >= 0 {
		var filtered []*model.Payment
		for _, p := range all {
			if p.Status == statusFilter {
				filtered = append(filtered, p)
			}
		}
		all = filtered
	}

	total := len(all)

	// Paginate
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	start := (page - 1) * limit
	if start >= total {
		return []*model.Payment{}, total, nil
	}
	end := start + limit
	if end > total {
		end = total
	}

	return all[start:end], total, nil
}

// RefundPayment marks an existing payment as refunded.
func (sp *StubProvider) RefundPayment(ctx context.Context, paymentID, reason string) (*model.Payment, error) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	// Retrieve payment
	p, err := sp.GetPayment(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("refund: %w", err)
	}

	if p.Status != model.Completed {
		return nil, fmt.Errorf(
			"cannot refund payment %s: current status is %s",
			paymentID, p.Status.Label(),
		)
	}

	// Simulate refund processing delay
	time.Sleep(time.Duration(50+sp.rng.Intn(100)) * time.Millisecond)

	p.Status = model.Refunded
	p.UpdatedAt = time.Now().UTC()

	// Persist updated status
	if sp.store != nil {
		if err := sp.store.UpdatePaymentStatus(ctx, paymentID, model.Refunded); err != nil {
			return nil, fmt.Errorf("persist refund: %w", err)
		}
	}
	sp.payments[paymentID] = p

	return p, nil
}
