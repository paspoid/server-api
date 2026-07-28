package use_cases

import (
	"context"

	"github.com/paspoid/server-api-go/internal/application/dto"
	"github.com/paspoid/server-api-go/internal/application/ports/providers"
	"github.com/paspoid/server-api-go/internal/application/ports/providers/requests"
)

type ValidateUseCase struct {
	transactionServiceProv providers.TransactionServiceProvider
}

func NewValidateUseCase(
	transactionServiceProv providers.TransactionServiceProvider,
) *ValidateUseCase {
	return &ValidateUseCase{
		transactionServiceProv: transactionServiceProv,
	}
}

func (uc *ValidateUseCase) Execute(
	ctx context.Context,
	inp dto.ValidateInput,
) (*dto.ValidateOutput, error) {
	validateResp, err := uc.transactionServiceProv.Validate(ctx, requests.ValidateRequest{
		Nonce: inp.Nonce,
	})
	if err != nil {
		return nil, err
	}

	return &dto.ValidateOutput{
		Status:     validateResp.Status,
		DataType:   validateResp.DataType,
		DataValue:  validateResp.DataValue,
		PhoneData:  validateResp.PhoneData,
		DeviceData: validateResp.DeviceData,
	}, nil
}
