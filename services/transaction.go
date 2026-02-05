package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	// helper.ExecuteTransaction(func() error)
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GenerateReport(fromDate *time.Time, toDate *time.Time) (*models.Report, error) {
	return s.repo.GenerateReport(fromDate, toDate)
}
