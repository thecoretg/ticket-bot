package connectwise

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func ValidateWebhook(r *http.Request) (bool, error) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return false, fmt.Errorf("reading request body: %w", err)
	}
	defer r.Body.Close()

	var meta struct {
		Metadata struct {
			KeyURL string `json:"key_url"`
		} `json:"metadata"`
	}
	if err := json.Unmarshal(payload, &meta); err != nil {
		return false, fmt.Errorf("unmarshaling request body: %w", err)
	}

	resp, err := http.Get(meta.Metadata.KeyURL)
	if err != nil {
		return false, fmt.Errorf("getting shared secret key from %s: %w", meta.Metadata.KeyURL, err)
	}
	defer resp.Body.Close()

	var keyResp struct {
		SigningKey string `json:"signing_key"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&keyResp); err != nil {
		return false, fmt.Errorf("decoding shared secret key response: %w", err)
	}

	sharedSecret := []byte(keyResp.SigningKey)
	hash := sha256.Sum256(sharedSecret)
	h := hmac.New(sha256.New, hash[:])
	h.Write(payload)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	expectedSig := r.Header.Get("x-content-signature")

	return signature == expectedSig, nil
}
