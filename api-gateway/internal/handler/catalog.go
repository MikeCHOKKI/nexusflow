// ─── Catalog Handlers ──────────────────────────────────────────
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/MikeCHOKKI/nexusflow/api-gateway/internal/client"
	pb "github.com/MikeCHOKKI/nexusflow/api-gateway/gen"
	"github.com/gorilla/mux"
)

// CatalogHandler exposes product CRUD endpoints.
type CatalogHandler struct {
	clients *client.Clients
}

func NewCatalogHandler(clients *client.Clients) *CatalogHandler {
	return &CatalogHandler{clients: clients}
}

// GET /api/v1/products
func (h *CatalogHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	req := &pb.ListProductsRequest{
		Page:     int32(parseInt(q.Get("page"), 1)),
		Limit:    int32(parseInt(q.Get("limit"), 20)),
		Category: q.Get("category"),
		Search:   q.Get("search"),
	}

	resp, err := h.clients.Catalog.ListProducts(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list products: "+err.Error())
		return
	}

	writeSuccess(w, http.StatusOK, resp)
}

// GET /api/v1/products/{id}
func (h *CatalogHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	product, err := h.clients.Catalog.GetProduct(r.Context(), &pb.GetProductRequest{Id: id})
	if err != nil {
		writeError(w, http.StatusNotFound, "product not found: "+err.Error())
		return
	}

	writeSuccess(w, http.StatusOK, product)
}

// POST /api/v1/products
func (h *CatalogHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req pb.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "product name is required")
		return
	}

	product, err := h.clients.Catalog.CreateProduct(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create product: "+err.Error())
		return
	}

	writeSuccess(w, http.StatusCreated, product)
}

// PUT /api/v1/products/{id}
func (h *CatalogHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req pb.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	req.Id = id

	product, err := h.clients.Catalog.UpdateProduct(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update product: "+err.Error())
		return
	}

	writeSuccess(w, http.StatusOK, product)
}

// DELETE /api/v1/products/{id}
func (h *CatalogHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	_, err := h.clients.Catalog.DeleteProduct(r.Context(), &pb.DeleteProductRequest{Id: id})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete product: "+err.Error())
		return
	}

	writeSuccess(w, http.StatusNoContent, nil)
}
