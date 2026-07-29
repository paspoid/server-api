package paspoid

import (
	"context"

	"github.com/paspoid/server-api-go/internal/application/dto"
	"github.com/paspoid/server-api-go/internal/application/use_cases"
	restProviders "github.com/paspoid/server-api-go/internal/infrastructure/adapters/providers/rest"
)

type Client struct {
	apiKey     string
	apiSecret  string
	getKeyUc   *use_cases.GetKeyUseCase
	validateUc *use_cases.ValidateUseCase
}

func NewClient(
	baseUrl,
	apiKey,
	apiSecret string,
) *Client {
	provider := restProviders.NewTransactionServiceRestProvider(
		baseUrl,
		apiKey,
		apiSecret,
	)

	return &Client{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		getKeyUc:   use_cases.NewGetKeyUseCase(provider),
		validateUc: use_cases.NewValidateUseCase(provider),
	}
}

func (c *Client) GetKey(
	ctx context.Context,
	servicePublicId,
	transactionType string,
) (*GetKeyResponse, error) {
	out, err := c.getKeyUc.Execute(ctx, dto.GetKeyInput{
		ApiKey:          c.apiKey,
		ApiSecret:       c.apiSecret,
		ServicePublicId: servicePublicId,
		TransactionType: transactionType,
	})
	if err != nil {
		return nil, err
	}

	return &GetKeyResponse{
		Key:              out.Key,
		ValidationWindow: out.ValidationWindow,
	}, nil
}

func (c *Client) Validate(
	ctx context.Context,
	nonce string,
) (*ValidateResponse, error) {
	out, err := c.validateUc.Execute(ctx, dto.ValidateInput{
		Nonce: nonce,
	})
	if err != nil {
		return nil, err
	}

	return &ValidateResponse{
		Status:     out.Status,
		DataType:   out.DataType,
		DataValue:  out.DataValue,
		PhoneData:  out.PhoneData,
		DeviceData: out.DeviceData,
	}, nil
}
