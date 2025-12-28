package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/celpung/gocleanarch/internal_old/delivery/handler"
)

func CompanyRouter(r chi.Router, companyHandler *handler.CompanyHandler) {
	r.Route("/companies", func(r chi.Router) {
		r.Post("/", companyHandler.CreateCompany)
		r.Get("/", companyHandler.ListCompanies)
		r.Get("/{id}", companyHandler.GetCompany)
		r.Put("/{id}", companyHandler.UpdateCompany)
		r.Delete("/{id}", companyHandler.DeleteCompany)
	})
}
