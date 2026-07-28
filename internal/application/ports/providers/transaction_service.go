package providers

import (
	"context"

	"github.com/paspoid/server-api-go/internal/application/ports/providers/requests"
	"github.com/paspoid/server-api-go/internal/application/ports/providers/responses"
)

type TransactionServiceProvider interface {
	GetKey(ctx context.Context, req requests.GetKeyRequest) (*responses.GetKeyResponse, error)
	Validate(ctx context.Context, req requests.ValidateRequest) (*responses.ValidateResponse, error)
}