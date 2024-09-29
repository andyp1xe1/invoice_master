package main

import (
	"encoding/json"
	"net/http"
	"gorm.io/gorm"
	"github.com/go-chi/chi/v5"
	 "strconv"
)

// InvoiceHandler handles the request to fetch all invoices
type InvoiceHandler struct {
	db *gorm.DB
}

// NewInvoiceHandler creates a new instance of InvoiceHandler
func NewInvoiceHandler(db *gorm.DB) *InvoiceHandler {
	return &InvoiceHandler{db: db}
}

// GetAllInvoices handles GET requests to retrieve all invoices in JSON
func (h *InvoiceHandler) GetAllInvoices(w http.ResponseWriter, r *http.Request) {
	var invoices []Contract
	// Query the database for all invoices
	if err := h.db.Preload("Services").Find(&invoices).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert the invoices slice to JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(invoices); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


// GetInvoiceByID handles GET requests to retrieve a specific invoice by its ID
func (h *InvoiceHandler) GetInvoiceByID(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL parameters
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	// Find the invoice by ID in the database
	var invoice Contract
	if err := h.db.Preload("Services").First(&invoice, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Invoice not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Convert the invoice to JSON and return
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(invoice); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


// GetAllContracts handles GET requests to retrieve all contracts in JSON
func (h *InvoiceHandler) GetAllContracts(w http.ResponseWriter, r *http.Request) {
	var contracts []Contract
	// Query the database for all contracts
	if err := h.db.Preload("Services").Find(&contracts).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert the contracts slice to JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(contracts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}