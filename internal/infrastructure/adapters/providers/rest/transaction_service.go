package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/paspoid/server-api/internal/application/ports/providers"
	"github.com/paspoid/server-api/internal/application/ports/providers/requests"
	"github.com/paspoid/server-api/internal/application/ports/providers/responses"
)

type transactionServiceRestProvider struct {
	baseUrl   string
	apiKey    string
	apiSecret string
	client    *http.Client
}

func NewTransactionServiceRestProvider(
	baseUrl string,
	apiKey string,
	apiSecret string,
) providers.TransactionServiceProvider {
	return &transactionServiceRestProvider{
		baseUrl:   strings.TrimRight(baseUrl, "/"),
		apiKey:    apiKey,
		apiSecret: apiSecret,
		client:    &http.Client{},
	}
}

func (p *transactionServiceRestProvider) GetKey(
	ctx context.Context,
	req requests.GetKeyRequest,
) (*responses.GetKeyResponse, error) {
	reqData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal get-key request: %w", err)
	}

	var payload map[string]any
	if err = json.Unmarshal(reqData, &payload); err != nil {
		return nil, fmt.Errorf("prepare get-key request: %w", err)
	}

	payload["api_key"] = p.apiKey
	payload["api_secret"] = p.apiSecret

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal get-key payload: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		p.baseUrl+"/v1/ext/get-key",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("create get-key request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("execute get-key request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read get-key response: %w", err)
	}

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf(
			"get-key returned status %d: %s",
			resp.StatusCode,
			strings.TrimSpace(string(respBody)),
		)
	}

	var result responses.GetKeyResponse
	if err = json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decode get-key response: %w", err)
	}

	return &result, nil
}

func (p *transactionServiceRestProvider) Validate(
	ctx context.Context,
	req requests.ValidateRequest,
) (*responses.ValidateResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal validate request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		p.baseUrl+"/v1/ext/validate",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("create validate request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("execute validate request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read validate response: %w", err)
	}

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf(
			"validate returned status %d: %s",
			resp.StatusCode,
			strings.TrimSpace(string(respBody)),
		)
	}

	var result responses.ValidateResponse
	if err = json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decode validate response: %w", err)
	}

	return &result, nil
}
