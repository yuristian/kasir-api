package routes

import (
	"database/sql"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"kasir-api/utils"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterAllRoutes(db *sql.DB) {
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)
	http.HandleFunc("/api/categories/{id}/produk", productHandler.GetAllProductsByCategoryID)

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)
	http.HandleFunc("/api/report/hari-ini", transactionHandler.GetReportToday)
	http.HandleFunc("/api/report", transactionHandler.GetReportByDate)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		healthCheck(w, r)
	})
}

// healthCheck godoc
// @Summary      Health Check
// @Tags         system
// @Router       /health [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "OK"})
}
