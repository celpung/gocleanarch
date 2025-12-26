package dto

type CreateCompanyRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type UpdateCompanyRequest struct {
	Name     *string `json:"name"`
	Address  *string `json:"address"`
	Phone    *string `json:"phone"`
	IsActive *bool   `json:"is_active"`
}

type CompanyResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	IsActive bool   `json:"is_active"`
}

type ListCompaniesResponse struct {
	Data []CompanyResponse `json:"data"`
	Meta PagingMeta        `json:"meta"`
}
