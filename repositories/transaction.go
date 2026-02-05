package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin() //untuk transaction
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)
	for _, item := range items {
		var productName string
		var productID, price, stok int
		err := tx.QueryRow("SELECT id, name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productID, &productName, &price, &stok)

		//return error kalau product not found
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}

		if err != nil {
			return nil, err
		}

		//hitung subtotal
		subtotal := item.Quantity * price
		totalAmount += subtotal

		//update stok
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		//masukin ke transaction details
		details = append(details, models.TransactionDetail{
			ProductID:   productID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING ID", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	if len(details) > 0 {
		query := "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES "
		values := []interface{}{}

		for i, d := range details {
			n := i * 4
			query += fmt.Sprintf("($%d, $%d, $%d, $%d),", n+1, n+2, n+3, n+4)
			values = append(values, transactionID, d.ProductID, d.Quantity, d.Subtotal)
			details[i].TransactionID = transactionID
		}

		query = query[:len(query)-1]
		_, err = tx.Exec(query, values...)
		// fmt.Printf(query)
		if err != nil {
			return nil, err
		}
	}

	// for i, detail := range details {
	// 	details[i].TransactionID = transactionID
	// 	_, err := tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
	// 		transactionID, detail.ProductID, detail.Quantity, detail.Subtotal)

	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	//Commit, and return error if commit failed
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) GenerateReport(fromDate *time.Time, toDate *time.Time) (*models.Report, error) {
	start := time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, fromDate.Location())
	var end time.Time
	if toDate == nil {
		end = time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 23, 59, 59, 0, fromDate.Location())
	} else {
		end = time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 23, 59, 59, 0, toDate.Location())
	}

	report := &models.Report{}
	querySummary := `
        SELECT 
            COALESCE(SUM(total_amount), 0), 
            COUNT(id)
        FROM transactions 
        WHERE created_at BETWEEN $1 AND $2`
	fmt.Println(querySummary, start, end)
	err := repo.db.QueryRow(querySummary, start, end).Scan(
		&report.TotalRevenue,
		&report.TotalTransaction,
	)
	if err != nil {
		return nil, err
	}

	queryPopular := `
        SELECT p.name, SUM(td.quantity) as total_sold
        	FROM transaction_details td
        	JOIN transactions t ON td.transaction_id = t.id
        	JOIN products p ON td.product_id = p.id
        	WHERE t.created_at BETWEEN $1 AND $2
        	GROUP BY p.id, p.name
        	ORDER BY total_sold DESC
        	LIMIT 1`

	err = repo.db.QueryRow(queryPopular, start, end).Scan(
		&report.PopularProduct.ProductName,
		&report.PopularProduct.Quantity,
	)

	// Jika tidak ada produk terjual sama sekali, QueryRow akan return sql.ErrNoRows
	if err == sql.ErrNoRows {
		report.PopularProduct = models.SoldProduct{ProductName: "N/A", Quantity: 0}
	} else if err != nil {
		return nil, err
	}

	return report, nil
}
