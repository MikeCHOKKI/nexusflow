// ─── Payment Handlers ──────────────────────────────────────────
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MikeCHOKKI/nexusflow/api-gateway/internal/client"
	"github.com/MikeCHOKKI/nexusflow/api-gateway/internal/middleware"
	pb "github.com/MikeCHOKKI/nexusflow/api-gateway/gen"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

// PaymentHandler exposes payment endpoints.
type PaymentHandler struct {
	clients *client.Clients
	rdb     *redis.Client
}

func NewPaymentHandler(clients *client.Clients, rdb *redis.Client) *PaymentHandler {
	return &PaymentHandler{clients: clients, rdb: rdb}
}

// POST /api/v1/payments
func (h *PaymentHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	userEmail := middleware.UserEmailFromContext(r.Context())

	var body struct {
		OrderID  string  `json:"order_id"`
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
		Method   string  `json:"method"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	method, ok := pb.PaymentMethod_value[body.Method]
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid payment method: "+body.Method)
		return
	}

	resp, err := h.clients.Payment.ProcessPayment(r.Context(), &pb.ProcessPaymentRequest{
		OrderId:  body.OrderID,
		UserId:   userID,
		Amount:   body.Amount,
		Currency: body.Currency,
		Method:   pb.PaymentMethod(method),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "payment failed: "+err.Error())
		return
	}

	// Publier un événement payment.received sur Redis pour le Notification Service
	if resp.GetPayment() != nil {
		p := resp.GetPayment()
		userName := middleware.UserNameFromContext(r.Context())
		eventData := map[string]interface{}{
			"order_id":        p.OrderId,
			"payment_id":      p.Id,
			"user_id":         userID,
			"customer_email":  userEmail,
			"customer_name":   userName,
			"amount":          p.Amount,
			"currency":        p.Currency,
			"payment_method":  body.Method,
			"transaction_id":  p.TransactionRef,
			"status":          p.Status.String(),
		}

		payload := map[string]interface{}{
			"type":      "payment.received",
			"data":      eventData,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}

		payloadJSON, _ := json.Marshal(payload)
		channel := "nexusflow:notifications"

		if err := h.rdb.Publish(context.Background(), channel, payloadJSON).Err(); err != nil {
			// Non-bloquant : on logge l'erreur mais on ne bloque pas la réponse
			fmt.Printf("[gateway] redis publish payment.received: %v\n", err)
		}
	}

	writeSuccess(w, http.StatusOK, resp)
}

// GET /api/v1/payments/{id}
func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	payment, err := h.clients.Payment.GetPayment(r.Context(), &pb.GetPaymentRequest{
		Id: id,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "payment not found: "+err.Error())
		return
	}

	writeSuccess(w, http.StatusOK, payment)
}

// POST /api/v1/payments/{id}/refund
func (h *PaymentHandler) RefundPayment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var body struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	resp, err := h.clients.Payment.RefundPayment(r.Context(), &pb.RefundPaymentRequest{
		PaymentId: id,
		Reason:    body.Reason,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "refund failed: "+err.Error())
		return
	}

	writeSuccess(w, http.StatusOK, resp)
}
