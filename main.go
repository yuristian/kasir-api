package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var produk = []Produk{
	{ID: 1, Nama: "Indomie Godog", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Vit 1000ml", Harga: 3000, Stok: 40},
	{ID: 3, Nama: "Kecap", Harga: 12000, Stok: 20},
}

var category = []Category{
	{ID: 1, Name: "Elektronik", Description: "Perangkat modern mulai dari smartphone hingga peralatan rumah tangga pintar."},
	{ID: 2, Name: "Pakaian Pria", Description: "Koleksi fashion pria termasuk kemeja, celana, dan aksesoris formal maupun kasual."},
	{ID: 3, Name: "Pakaian Wanita", Description: "Tren busana wanita terbaru, gaun, atasan, dan perlengkapan fashion lainnya."},
	{ID: 4, Name: "Kesehatan & Kecantikan", Description: "Produk perawatan kulit, kosmetik, dan suplemen kesehatan tubuh."},
	{ID: 5, Name: "Peralatan Rumah Tangga", Description: "Kebutuhan dapur, dekorasi ruang tamu, dan alat kebersihan rumah."},
	{ID: 6, Name: "Olahraga & Outdoor", Description: "Alat bantu olahraga, pakaian atletik, dan perlengkapan berkemah."},
	{ID: 7, Name: "Otomotif", Description: "Suku cadang kendaraan, aksesoris mobil, dan alat perawatan mesin."},
	{ID: 8, Name: "Buku & Alat Tulis", Description: "Berbagai genre buku bacaan serta perlengkapan kantor dan sekolah."},
	{ID: 9, Name: "Mainan & Hobi", Description: "Koleksi mainan anak-anak, puzzle, dan perlengkapan hobi kreatif."},
	{ID: 10, Name: "Makanan & Minuman", Description: "Produk kuliner, camilan, dan bahan pangan segar maupun instan."},
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// GET localhost:8080/api/produk/{id}
func getProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid Produk ID")
		return
	}

	for _, p := range produk {
		if p.ID == id {
			respondWithJSON(w, http.StatusOK, p)
			return
		}
	}

	respondWithError(w, http.StatusNotFound, "Produk Belum Tersedia")
}

// PUT localhost:8080/api/produk/{id}
func updateProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Produk ID")
		return
	}

	var updateProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	for i := range produk {
		if produk[i].ID == id {
			updateProduk.ID = id
			produk[i] = updateProduk

			respondWithJSON(w, http.StatusOK, updateProduk)
			return
		}
	}

	respondWithError(w, http.StatusNotFound, "Produk Belum Tersedia")
}

// DELETE localhost:8080/api/produk/{id}
func deleteProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Produk ID")
		return
	}

	for i, p := range produk {
		if p.ID == id {
			produk = append(produk[:i], produk[i+1:]...)
			respondWithJSON(w, http.StatusOK, "Delete Sukses")
			return
		}
	}

	respondWithError(w, http.StatusNotFound, "Produk Belum Tersedia")
}

// POST localhost:8080/api/produk
func insertNewProduk(w http.ResponseWriter, r *http.Request) {
	var produkBaru Produk
	err := json.NewDecoder(r.Body).Decode(&produkBaru)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	if len(produk) > 0 {
		produkBaru.ID = produk[len(produk)-1].ID + 1
	} else {
		produkBaru.ID = 1
	}

	produk = append(produk, produkBaru)
	respondWithJSON(w, http.StatusCreated, produkBaru)
}

// POST localhost:8080/api/category
func insertNewCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	if len(category) > 0 {
		newCategory.ID = category[len(category)-1].ID + 1
	} else {
		newCategory.ID = 1
	}

	category = append(category, newCategory)
	respondWithJSON(w, http.StatusCreated, newCategory)
}

// PUT localhost:8080/api/category/{id}
func updateCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Category ID")
		return
	}

	var updateCategory Category
	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	for i := range category {
		if category[i].ID == id {
			updateCategory.ID = id
			category[i] = updateCategory
			respondWithJSON(w, http.StatusOK, updateCategory)
			return
		}
	}

	respondWithError(w, http.StatusNotFound, "Invalid Category ID")
}

// GET localhost:8080/api/category/{id}
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Category ID")
		return
	}

	for _, cat := range category {
		if cat.ID == id {
			respondWithJSON(w, http.StatusOK, cat)
			return
		}
	}

	respondWithError(w, http.StatusNotFound, "Invalid Category ID")
}

// DELETE localhost:8080/api/category/{id}
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Category ID")
		return
	}

	for i, cat := range category {
		if cat.ID == id {
			category = append(category[:i], category[i+1:]...)
			respondWithJSON(w, http.StatusOK, cat)
			return
		}
	}

	respondWithError(w, http.StatusNotFound, "Invalid Category ID")
}

func main() {
	// GET localhost:8080/api/produk/{id}
	// PUT localhost:8080/api/produk/{id}
	// DELETE localhost:8080/api/produk/{id}
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProdukByID(w, r)
		} else if r.Method == "PUT" {
			updateProdukByID(w, r)
		} else if r.Method == "DELETE" {
			deleteProduk(w, r)
		}
	})

	// GET localhost:8080/api/produk
	// POST localhost:8080/api/produk
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
		} else if r.Method == "POST" {
			insertNewProduk(w, r)
		}
	})

	//GET /api/categories
	//POST /api/categories
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category)
		} else if r.Method == "POST" {
			insertNewCategory(w, r)
		}
	})

	//PUT /api/categories/:id
	//GET /api/categories/:id
	//DELETE /api/categories/:id
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			updateCategoryByID(w, r)
		} else if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}
	})

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API running",
		})
	})

	fmt.Println("server running di localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
