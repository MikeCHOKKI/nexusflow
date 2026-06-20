// ─── Order Handlers ────────────────────────────────────────────
package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MikeCHOKKI/nexusflow/api-gateway/internal/middleware"
	"github.com/gorilla/mux"
)

// unwrapPHPResponse extracts the "data" field from a PHP service response
// (which wraps everything in {"success":true, "data": ...}) to avoid double-wrapping
// when the gateway adds its own envelope via writeJSON/writeSuccess.
func unwrapPHPResponse(result map[string]interface{}) interface{} {
	if data, ok := result["data"]; ok {
		return data
	}
	return result
}

// OrderHandler exposes order endpoints over REST (Order Service is PHP/HTTP).
type OrderHandler struct {
	orderAddr string
}

func NewOrderHandler(orderAddr string) *OrderHandler {
	return &OrderHandler{orderAddr: orderAddr}
}

func (h *OrderHandler) restURL(path string) string {
	return fmt.Sprintf("http://%s%s", h.orderAddr, path)
}

// POST /api/v1/orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	userEmail := middleware.UserEmailFromContext(r.Context())

	body := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	body["user_id"] = userID
	body["customer_email"] = userEmail
	// customer_name optionnel — le mail l'utilise s'il est fourni
	if name, ok := body["customer_name"]; !ok || name == "" {
		body["customer_name"] = "Client"
	}

	data, _ := json.Marshal(body)
	resp, err := http.Post(h.restURL("/orders"), "application/json", bytes.NewReader(data))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create order: "+err.Error())
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if resp.StatusCode >= 400 {
		writeError(w, resp.StatusCode, fmt.Sprintf("order service error: %v", result))
		return
	}
	writeSuccess(w, resp.StatusCode, unwrapPHPResponse(result))
}

// GET /api/v1/orders
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	url := h.restURL("/orders?user_id=" + userID)
	if p := r.URL.Query().Get("page"); p != "" {
		url += "&page=" + p
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		url += "&limit=" + l
	}

	resp, err := http.Get(url)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list orders: "+err.Error())
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	writeSuccess(w, resp.StatusCode, unwrapPHPResponse(result))
}

// GET /api/v1/orders/{id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userID := middleware.UserIDFromContext(r.Context())
	url := h.restURL("/orders/" + id + "?user_id=" + userID)

	resp, err := http.Get(url)
	if err != nil {
		writeError(w, http.StatusNotFound, "order not found: "+err.Error())
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	writeSuccess(w, resp.StatusCode, unwrapPHPResponse(result))
}

// PUT /api/v1/orders/{id}/status
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, h.restURL("/orders/"+id+"/status"), bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update order status: "+err.Error())
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	writeSuccess(w, resp.StatusCode, unwrapPHPResponse(result))
}
