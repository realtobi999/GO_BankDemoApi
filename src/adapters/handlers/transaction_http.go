package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
	"github.com/realtobi999/GO_BankDemoApi/src/core/ports"
)

type TransactionHandler struct {
	TransactionService ports.ITransactionService
}

func NewTransactionHandler(transactionService ports.ITransactionService) *TransactionHandler {
	return &TransactionHandler{
		TransactionService: transactionService,
	}
}

func (h *TransactionHandler) Index(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parseLimitOffsetParams(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse parameters: "+err.Error())
		return
	}

	accountID := uuid.Nil
	if r.URL.Query().Get("account_id") != ""{
		accountID, err = uuid.Parse(r.URL.Query().Get("account_id"))
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
			return
		}
	}

	transactions, err := h.TransactionService.Index(accountID, limit, offset)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		} else if errors.Is(err, domain.ErrInternalFailure) {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

	}

	RespondWithJsonAndSerializeList(w, http.StatusOK, transactions)
}

func (h *TransactionHandler) Get(w http.ResponseWriter, r *http.Request) {
	transactionID, err := uuid.Parse(chi.URLParam(r, "transaction_id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
		return
	}

	transaction, err := h.TransactionService.Get(transactionID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, domain.ErrInternalFailure) {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	RespondWithJsonAndSerialize(w, http.StatusOK, transaction)
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse UUID: "+err.Error())
		return
	}

	body, err := decode[domain.CreateTransactionRequest](r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse the body: "+err.Error())
		return
	}
	body.SenderAccountID = accountID

	transaction, err := h.TransactionService.Create(body)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, domain.ErrBadRequest) {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, domain.ErrValidation) {
			RespondWithValidationErrors(w, http.StatusBadRequest, "Failed to validate request", domain.ExtractValidationErrorsToList(err))
			return
		}		
		if errors.Is(err, domain.ErrInternalFailure) {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	w.Header().Set("Location", fmt.Sprintf("/api/transaction/%s", transaction.ID.String()))
	RespondWithJson(w, http.StatusCreated, nil)
}