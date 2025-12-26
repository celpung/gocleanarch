package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/celpung/gocleanarch/internal/delivery/dto"
	"github.com/celpung/gocleanarch/internal/usecase/port"
	"github.com/celpung/gocleanarch/pkg/helpers/httpx"
)

type CompanyHandler struct {
	usecase port.CompanyUsecase
}

func NewCompanyHandler(usecase port.CompanyUsecase) *CompanyHandler {
	return &CompanyHandler{usecase: usecase}
}

func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCompanyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request body"})
		return
	}

	input := port.CreateCompanyInput{
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
	}

	company, err := h.usecase.CreateCompany(r.Context(), input)
	if err != nil {
		respondCompanyError(w, err)
		return
	}

	resp := dto.CompanyResponse{
		ID:       company.ID,
		Name:     company.Name,
		Address:  company.Address,
		Phone:    company.Phone,
		IsActive: company.IsActive,
	}

	httpx.WriteJSON(w, http.StatusCreated, resp)
}

func (h *CompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	company, err := h.usecase.GetCompany(r.Context(), id)
	if err != nil {
		respondCompanyError(w, err)
		return
	}

	resp := dto.CompanyResponse{
		ID:       company.ID,
		Name:     company.Name,
		Address:  company.Address,
		Phone:    company.Phone,
		IsActive: company.IsActive,
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *CompanyHandler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	page, limit := 1, 10

	if v := r.URL.Query().Get("page"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			page = parsed
		}
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			limit = parsed
		}
	}

	companies, total, err := h.usecase.ListCompanies(r.Context(), page, limit)
	if err != nil {
		respondCompanyError(w, err)
		return
	}

	items := make([]dto.CompanyResponse, 0, len(companies))
	for _, c := range companies {
		items = append(items, dto.CompanyResponse{
			ID:       c.ID,
			Name:     c.Name,
			Address:  c.Address,
			Phone:    c.Phone,
			IsActive: c.IsActive,
		})
	}

	resp := dto.ListCompaniesResponse{
		Data: items,
		Meta: dto.PagingMeta{
			Page:       page,
			Limit:      limit,
			TotalItems: total,
		},
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *CompanyHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.UpdateCompanyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request body"})
		return
	}

	input := port.UpdateCompanyInput{
		Name:     req.Name,
		Address:  req.Address,
		Phone:    req.Phone,
		IsActive: req.IsActive,
	}

	if err := h.usecase.UpdateCompany(r.Context(), id, input); err != nil {
		respondCompanyError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *CompanyHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.usecase.DeleteCompany(r.Context(), id); err != nil {
		respondCompanyError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusNoContent, nil)
}

func respondCompanyError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	resp := map[string]string{
		"message": "internal server error",
	}

	var vErr port.ValidationError

	switch {
	case errors.As(err, &vErr):
		status = http.StatusBadRequest
		resp["message"] = vErr.Error()
	case errors.Is(err, port.ErrCompanyNotFound):
		status = http.StatusNotFound
		resp["message"] = "company not found"
	case errors.Is(err, port.ErrCompanyNameExists):
		status = http.StatusConflict
		resp["message"] = "company name already exists"
	case errors.Is(err, port.ErrNoFieldsToUpdate):
		status = http.StatusBadRequest
		resp["message"] = "no fields to update"
	default:
	}

	httpx.WriteJSON(w, status, resp)
}
