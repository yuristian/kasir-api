package handlers

import (
	"encoding/json"
	"fmt"
	"kasir-api/models"
	"kasir-api/services"
	"kasir-api/utils"
	"net/http"
	"time"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *TransactionHandler) GetReportToday(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GenerateTodayReport(w, r)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *TransactionHandler) GetReportByDate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GenerateReportByDate(w, r)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// Checkout godoc
// @Summary      Proses Checkout Transaksi
// @Tags         Transaction
// @Accept       json
// @Produce      json
// @Param        payload body models.CheckoutRequest true "Payload Checkout"
// @Success      201  {object}  models.Transaction
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/checkout [post]
func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req models.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	fmt.Print(err)
	if err != nil {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	tx, err := h.service.Checkout(req.Items)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, tx)
}

func (h *TransactionHandler) GenerateTodayReport(w http.ResponseWriter, r *http.Request) {
	date := time.Now()
	report, err := h.service.GenerateReport(&date, nil)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, report)
}

func (h *TransactionHandler) GenerateReportByDate(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start_date")
	endStr := r.URL.Query().Get("end_date")

	if startStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "start_date is required")
		return
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid start_date format (YYYY-MM-DD)")
		return
	}

	var endPtr *time.Time
	if endStr != "" {
		parsedEnd, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid end_date format")
			return
		}
		endPtr = &parsedEnd
	}

	fmt.Println(&start, endPtr)

	report, err := h.service.GenerateReport(&start, endPtr)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, report)
}
