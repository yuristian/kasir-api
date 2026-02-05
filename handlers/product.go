package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"kasir-api/utils"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAllProducts(w, r)
	case http.MethodPost:
		h.CreateProduct(w, r)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetProductByID(w, r)
	case http.MethodPut:
		h.UpdateProduct(w, r)
	case http.MethodDelete:
		h.DeleteProduct(w, r)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// GetAllProducts godoc
// @Summary      Daftar Semua Produk
// @Tags         product
// @Produce      json
// @Success      200  {array}  models.Product
// @Router       /api/product [get]
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	products, err := h.service.GetAllProducts(name)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, products)
}

// GetAllProductsByCategoryID godoc
// @Summary      Daftar Semua Produk berdasarkan kategori ID
// @Tags         product
// @Produce      json
// @Success      200  {array}  models.Product
// @Router       /api/categories/{id}/produk [get]
func (h *ProductHandler) GetAllProductsByCategoryID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	products, err := h.service.GetAllProductsByCategoryID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Produk tidak ditemukan")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, products)
}

// CreateProduct godoc
// @Summary      Tambah Produk Baru
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        product  body      models.Product  true  "Data Product"
// @Success      201     {object}  models.Product
// @Router       /api/product [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	product := models.Product{}
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.CreateProduct(&product)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, product)
}

// getProductByID godoc
// @Summary      Ambil Produk by ID
// @Tags         product
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  models.Product
// @Router       /api/product/{id} [get]
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.service.GetProductByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, product)
}

// updateProduct godoc
// @Summary      Update Produk
// @Tags         product
// @Param        id    path      int     true  "Product ID"
// @Param        data  body      models.Product  true  "Data Update"
// @Router       /api/product/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product.ID = id
	err = h.service.UpdateProduct(&product)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, product)
}

// deleteProduct godoc
// @Summary      Hapus Produk
// @Tags         product
// @Param        id   path      int  true  "Product ID"
// @Router       /api/product/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	err = h.service.DeleteProduct(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Product deleted successfully",
	})
}
