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

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service}
}

func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAllCategories(w, r)
	case http.MethodPost:
		h.CreateCategory(w, r)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetCategoryByID(w, r)
	case http.MethodPut:
		h.UpdateCategory(w, r)
	case http.MethodDelete:
		h.DeleteCategory(w, r)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// GetAllCategories godoc
// @Summary      Daftar Semua Kategori
// @Tags         categories
// @Router       /api/categories [get]
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, categories)
}

// CreateCategory godoc
// @Summary      Tambah Kategori
// @Tags         categories
// @Param        data  body  models.Category  true  "Data Kategori"
// @Router       /api/categories [post]
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	category := models.Category{}
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.CreateCategory(&category)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, category)
}

// getCategoryByID godoc
// @Summary      Ambil Kategori by ID
// @Tags         categories
// @Param        id  path  int  true  "Kategori ID"
// @Success      200  {object}  models.Category
// @Router       /api/categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, category)
}

// updateCategoryByID godoc
// @Summary      Update Kategori
// @Tags         categories
// @Param        id  path  int  true  "Kategori ID"
// @Router       /api/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var category models.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	category.ID = id
	err = h.service.UpdateCategory(&category)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, category)
}

// deleteCategory godoc
// @Summary      Hapus Kategori
// @Tags         categories
// @Param        id  path  int  true  "Kategori ID"
// @Router       /api/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	err = h.service.DeleteCategory(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Category deleted successfully",
	})
}
