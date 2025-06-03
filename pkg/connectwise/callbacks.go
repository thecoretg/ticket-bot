package connectwise

import (
	"context"
	"fmt"
)

func callbackIdEndpoint(boardId int) string {
	return fmt.Sprintf("system/callbacks/%d", boardId)
}

func (c *Client) ListCallbacks(ctx context.Context, params *QueryParams) ([]Callback, error) {
	return ApiRequestPaginated[Callback](ctx, c, "GET", "system/callbacks", params, nil)
}

func (c *Client) PostCallback(ctx context.Context, callback *Callback) (*Callback, error) {
	return ApiRequestNonPaginated[Callback](ctx, c, "POST", "system/callbacks", nil, callback)
}

func (c *Client) GetCallback(ctx context.Context, callbackId int, params *QueryParams) (*Callback, error) {
	return ApiRequestNonPaginated[Callback](ctx, c, "GET", callbackIdEndpoint(callbackId), params, nil)
}

func (c *Client) DeleteCallback(ctx context.Context, callbackId int) error {
	if _, err := ApiRequestNonPaginated[struct{}](ctx, c, "DELETE", callbackIdEndpoint(callbackId), nil, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) PutCallback(ctx context.Context, callbackId int, callback *Callback) (*Callback, error) {
	return ApiRequestNonPaginated[Callback](ctx, c, "PUT", callbackIdEndpoint(callbackId), nil, callback)
}

func (c *Client) PatchCallback(ctx context.Context, callbackId int, patchOps []PatchOp) (*Callback, error) {
	return ApiRequestNonPaginated[Callback](ctx, c, "PATCH", callbackIdEndpoint(callbackId), nil, patchOps)
}
