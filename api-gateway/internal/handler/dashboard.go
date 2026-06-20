// ─── Dashboard Stats Handler ─────────────────────────────────────
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/MikeCHOKKI/nexusflow/api-gateway/internal/client"
	pb "github.com/MikeCHOKKI/nexusflow/api-gateway/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// DashboardHandler exposes aggregated stats for the frontend.
type DashboardHandler struct {
	clients  *client.Clients
	conns    []*grpc.ClientConn
	orderURL string
}

// NewDashboardHandler creates a handler that queries all services.
func NewDashboardHandler(clients *client.Clients, conns []*grpc.ClientConn, orderURL string) *DashboardHandler {
	return &DashboardHandler{clients: clients, conns: conns, orderURL: orderURL}
}

// ─── Types ──────────────────────────────────────────────────────

type dashboardStats struct {
	TotalProducts int               `json:"totalProducts"`
	TotalOrders   int               `json:"totalOrders"`
	TotalRevenue  float64           `json:"totalRevenue"`
	TotalUsers    int               `json:"totalUsers"`
	OrdersTrend   []trendPoint      `json:"ordersTrend"`
	RecentOrders  []orderSummary    `json:"recentOrders"`
	Services      []serviceHealth   `json:"services"`
}

type trendPoint struct {
	Date    string  `json:"date"`
	Count   int     `json:"count"`
	Revenue float64 `json:"revenue"`
}

type orderSummary struct {
	ID              string        `json:"id"`
	UserID          string        `json:"userId"`
	UserName        string        `json:"userName"`
	Items           []interface{} `json:"items"`
	Total           float64       `json:"total"`
	Status          string        `json:"status"`
	PaymentMethod   string        `json:"paymentMethod"`
	ShippingAddress string        `json:"shippingAddress"`
	CreatedAt       string        `json:"createdAt"`
	UpdatedAt       string        `json:"updatedAt"`
}

type serviceHealth struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Uptime  string `json:"uptime"`
	Latency int    `json:"latency"`
}

// ─── Stats ──────────────────────────────────────────────────────

// Stats responds to GET /api/v1/dashboard/stats.
func (h *DashboardHandler) Stats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	stats := h.gatherStats(ctx)
	writeJSON(w, http.StatusOK, stats)
}

func (h *DashboardHandler) gatherStats(ctx context.Context) dashboardStats {
	var (
		userCount    int
		productCount int
		orders       []map[string]interface{}
		revenue      float64
	)

	// ── User count via gRPC ─────────────────────────────────────
	usersResp, err := h.clients.User.ListUsers(ctx, &pb.ListUsersRequest{Limit: 1})
	if err != nil {
		log.Printf("[dashboard] ListUsers error: %v", err)
	} else if usersResp != nil && usersResp.Pagination != nil {
		userCount = int(usersResp.Pagination.Total)
	}

	// ── Product count via gRPC ──────────────────────────────────
	prodResp, err := h.clients.Catalog.ListProducts(ctx, &pb.ListProductsRequest{Limit: 1})
	if err != nil {
		log.Printf("[dashboard] ListProducts error: %v", err)
	} else if prodResp != nil && prodResp.Pagination != nil {
		productCount = int(prodResp.Pagination.Total)
	}

	// ── Orders via REST ─────────────────────────────────────────
	orders = h.fetchOrders(ctx)
	revenue = 0.0
	for _, o := range orders {
		if amt, ok := o["total_amount"].(float64); ok {
			revenue += amt
		}
	}

	recent := h.buildRecentOrders(orders)
	trend := h.buildTrend(orders)

	return dashboardStats{
		TotalProducts: productCount,
		TotalOrders:   len(orders),
		TotalRevenue:  revenue,
		TotalUsers:    userCount,
		OrdersTrend:   trend,
		RecentOrders:  recent,
		Services:      h.serviceHealth(),
	}
}

// ── Orders via REST ─────────────────────────────────────────────

type orderListResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Orders     []map[string]interface{} `json:"orders"`
		Pagination map[string]interface{}   `json:"pagination"`
	} `json:"data"`
}

func (h *DashboardHandler) fetchOrders(ctx context.Context) []map[string]interface{} {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://%s/orders?limit=100", h.orderURL), nil)
	if err != nil {
		log.Printf("[dashboard] fetch orders req: %v", err)
		return nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[dashboard] fetch orders: %v", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var parsed orderListResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		log.Printf("[dashboard] parse orders: %v", err)
		return nil
	}

	if !parsed.Success {
		return nil
	}

	return parsed.Data.Orders
}

// ── Recent orders ───────────────────────────────────────────────

func (h *DashboardHandler) buildRecentOrders(orders []map[string]interface{}) []orderSummary {
	// Sort by created_at desc
	sort.Slice(orders, func(i, j int) bool {
		ti, _ := orders[i]["created_at"].(string)
		tj, _ := orders[j]["created_at"].(string)
		return ti > tj
	})

	// Take last 5
	if len(orders) > 5 {
		orders = orders[:5]
	}

	result := make([]orderSummary, 0, len(orders))
	for _, o := range orders {
		result = append(result, orderSummary{
			ID:              strVal(o, "id"),
			UserID:          strVal(o, "user_id"),
			UserName:        strVal(o, "user_id")[:8], // short ID as name fallback
			Items:           []interface{}{},
			Total:           floatVal(o, "total_amount"),
			Status:          strVal(o, "status"),
			PaymentMethod:   "wave",
			ShippingAddress: strVal(o, "shipping_address"),
			CreatedAt:       strVal(o, "created_at"),
			UpdatedAt:       strVal(o, "updated_at"),
		})
	}
	return result
}

// ── Trend ───────────────────────────────────────────────────────

func (h *DashboardHandler) buildTrend(orders []map[string]interface{}) []trendPoint {
	// Aggregate by day (last 7 days)
	dayBuckets := make(map[string]*trendPoint)
	now := time.Now()

	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		label := now.AddDate(0, 0, -i).Format("Mon")
		dayBuckets[date] = &trendPoint{Date: label, Count: 0, Revenue: 0}
	}

	for _, o := range orders {
		created := parseTime(strVal(o, "created_at"))
		key := created.Format("2006-01-02")
		if bucket, ok := dayBuckets[key]; ok {
			bucket.Count++
			bucket.Revenue += floatVal(o, "total_amount")
		}
	}

	result := make([]trendPoint, 0, 7)
	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		label := now.AddDate(0, 0, -i).Format("Mon")
		if b, ok := dayBuckets[date]; ok {
			result = append(result, *b)
		} else {
			result = append(result, trendPoint{Date: label, Count: 0, Revenue: 0})
		}
	}
	return result
}

// ── Service health ──────────────────────────────────────────────

func (h *DashboardHandler) serviceHealth() []serviceHealth {
	serviceNames := []string{
		"api-gateway", "service-user", "service-catalog",
		"service-order", "service-payment", "service-notify",
	}

	result := make([]serviceHealth, 0, len(serviceNames))
	for _, name := range serviceNames {
		result = append(result, serviceHealth{
			Name:    name,
			Status:  "healthy",
			Uptime:  "99.9%",
			Latency: 5,
		})
	}

	// Update gRPC connection based health
	grpcServices := []struct {
		name string
		conn *grpc.ClientConn
	}{
		{"service-user", h.conns[0]},
		{"service-catalog", h.conns[1]},
		{"service-order", h.conns[2]},
		{"service-payment", h.conns[3]},
	}

	for i, s := range grpcServices {
		state := "healthy"
		latency := 5
		switch s.conn.GetState() {
		case connectivity.Ready:
			state = "healthy"
			latency = 8
		case connectivity.Connecting:
			state = "degraded"
			latency = 25
		case connectivity.TransientFailure:
			state = "down"
			latency = 0
		default:
			state = "degraded"
			latency = 15
		}
		for j := range result {
			if result[j].Name == s.name {
				result[j].Status = state
				result[j].Latency = latency
				break
			}
		}
		_ = i
	}

	return result
}

// ── Helpers ─────────────────────────────────────────────────────

func strVal(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	b, _ := json.Marshal(m[key])
	return string(bytes.Trim(b, `"`))
}

func floatVal(m map[string]interface{}, key string) float64 {
	switch v := m[key].(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		var f float64
		fmt.Sscanf(v, "%f", &f)
		return f
	}
	return 0
}

func parseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		return t
	}
	t, err = time.Parse("2006-01-02T15:04:05.999999-07:00", s)
	if err == nil {
		return t
	}
	t, err = time.Parse("2006-01-02 15:04:05.999999-07:00", s)
	if err == nil {
		return t
	}
	t, err = time.Parse("2006-01-02 15:04:05.999999-07:00", s)
	if err == nil {
		return t
	}
	return time.Now()
}
