package dto

type GetKeyInput struct {
	ApiKey          string
	ApiSecret       string
	ServicePublicId string
	TransactionType string
}

type GetKeyOutput struct {
	Key              string
	ValidationWindow string
}