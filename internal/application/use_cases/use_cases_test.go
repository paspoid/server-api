package use_cases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/paspoid/server-api-go/internal/application/dto"
	"github.com/paspoid/server-api-go/internal/application/ports/providers/requests"
	"github.com/paspoid/server-api-go/internal/application/ports/providers/responses"
	"github.com/paspoid/server-api-go/internal/application/use_cases"
)

type mockProvider struct {
	getKeyResp   *responses.GetKeyResponse
	getKeyErr    error
	validateResp *responses.ValidateResponse
	validateErr  error
}

func (m *mockProvider) GetKey(ctx context.Context, req requests.GetKeyRequest) (*responses.GetKeyResponse, error) {
	if m.getKeyErr != nil {
		return nil, m.getKeyErr
	}
	return m.getKeyResp, nil
}

func (m *mockProvider) Validate(ctx context.Context, req requests.ValidateRequest) (*responses.ValidateResponse, error) {
	if m.validateErr != nil {
		return nil, m.validateErr
	}
	return m.validateResp, nil
}

func TestGetKeyUseCase_Success(t *testing.T) {
	mock := &mockProvider{
		getKeyResp: &responses.GetKeyResponse{
			Key:              "key-123",
			ValidationWindow: "30s",
		},
	}

	uc := use_cases.NewGetKeyUseCase(mock)
	out, err := uc.Execute(context.Background(), dto.GetKeyInput{
		ServicePublicId: "srv-1",
		TransactionType: "phones",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.Key != "key-123" {
		t.Errorf("expected key 'key-123', got '%s'", out.Key)
	}
}

func TestGetKeyUseCase_Error(t *testing.T) {
	mock := &mockProvider{
		getKeyErr: errors.New("provider failure"),
	}

	uc := use_cases.NewGetKeyUseCase(mock)
	_, err := uc.Execute(context.Background(), dto.GetKeyInput{
		ServicePublicId: "srv-1",
		TransactionType: "phones",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestValidateUseCase_Success(t *testing.T) {
	mock := &mockProvider{
		validateResp: &responses.ValidateResponse{
			Status: "success",
		},
	}

	uc := use_cases.NewValidateUseCase(mock)
	out, err := uc.Execute(context.Background(), dto.ValidateInput{
		Nonce: "nonce-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.Status != "success" {
		t.Errorf("expected status 'success', got '%s'", out.Status)
	}
}
