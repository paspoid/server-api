package use_cases

import (
	"context"

	"github.com/paspoid/server-api/internal/application/dto"
	"github.com/paspoid/server-api/internal/application/ports/providers"
	"github.com/paspoid/server-api/internal/application/ports/providers/requests"
)

type GetKeyUseCase struct {
	transactionServiceProv providers.TransactionServiceProvider
}

func NewGetKeyUseCase(
	transactionServiceProv providers.TransactionServiceProvider,
) *GetKeyUseCase {
	return &GetKeyUseCase{
		transactionServiceProv: transactionServiceProv,
	}
}

func (uc *GetKeyUseCase) Execute(
	ctx context.Context,
	inp dto.GetKeyInput,
) (*dto.GetKeyOutput, error) {
	getKeyResp, err := uc.transactionServiceProv.GetKey(ctx, requests.GetKeyRequest{
		ServicePublicId: inp.ServicePublicId,
		TransactionType: inp.TransactionType,
	})
	if err != nil {
		return nil, err
	}

	return &dto.GetKeyOutput{
		Key:              getKeyResp.Key,
		ValidationWindow: getKeyResp.ValidationWindow,
	}, nil
}
