package connectwise

import (
	"context"
	"fmt"
)

func companyIdEndpoint(companyId int) string {
	return fmt.Sprintf("company/companies/%d", companyId)
}

func (c *Client) ListCompanies(ctx context.Context, params *QueryParams) ([]Company, error) {
	return ApiRequestPaginated[Company](ctx, c, "GET", "company/companies", params, nil)
}

func (c *Client) PostCompany(ctx context.Context, company *Company) (*Company, error) {
	return ApiRequestNonPaginated[Company](ctx, c, "POST", "company/companies", nil, company)
}

func (c *Client) GetCompany(ctx context.Context, companyId int, params *QueryParams) (*Company, error) {
	return ApiRequestNonPaginated[Company](ctx, c, "GET", companyIdEndpoint(companyId), params, nil)
}

func (c *Client) DeleteCompany(ctx context.Context, companyId int) error {
	if _, err := ApiRequestNonPaginated[struct{}](ctx, c, "DELETE", companyIdEndpoint(companyId), nil, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) PutCompany(ctx context.Context, companyId int, company *Company) (*Company, error) {
	return ApiRequestNonPaginated[Company](ctx, c, "PUT", companyIdEndpoint(companyId), nil, company)
}

func (c *Client) PatchCompany(ctx context.Context, companyId int, patchOps []PatchOp) (*Company, error) {
	return ApiRequestNonPaginated[Company](ctx, c, "PATCH", companyIdEndpoint(companyId), nil, patchOps)
}
