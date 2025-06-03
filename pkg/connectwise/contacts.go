package connectwise

import (
	"context"
	"fmt"
)

func contactIdEndpoint(contactId int) string {
	return fmt.Sprintf("company/contacts/%d", contactId)
}

func (c *Client) ListContacts(ctx context.Context, params *QueryParams) ([]Contact, error) {
	return ApiRequestPaginated[Contact](ctx, c, "GET", "company/contacts", params, nil)
}

func (c *Client) PostContact(ctx context.Context, contact *Contact) (*Contact, error) {
	return ApiRequestNonPaginated[Contact](ctx, c, "POST", "company/contact", nil, contact)
}

func (c *Client) GetContact(ctx context.Context, contactId int, params *QueryParams) (*Contact, error) {
	return ApiRequestNonPaginated[Contact](ctx, c, "GET", contactIdEndpoint(contactId), params, nil)
}

func (c *Client) DeleteContact(ctx context.Context, contactId int) error {
	if _, err := ApiRequestNonPaginated[struct{}](ctx, c, "DELETE", contactIdEndpoint(contactId), nil, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) PutContact(ctx context.Context, contactId int, contact *Contact) (*Contact, error) {
	return ApiRequestNonPaginated[Contact](ctx, c, "PUT", contactIdEndpoint(contactId), nil, contact)
}

func (c *Client) PatchContact(ctx context.Context, contactId int, patchOps []PatchOp) (*Contact, error) {
	return ApiRequestNonPaginated[Contact](ctx, c, "PATCH", contactIdEndpoint(contactId), nil, patchOps)
}
