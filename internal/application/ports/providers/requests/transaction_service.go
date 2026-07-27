package requests

type GetKeyRequest struct {
    ServicePublicId string `json:"service_public_id"`
    TransactionType string `json:"transaction_type"`
}

type ValidateRequest struct {
    Nonce string `json:"nonce"`
}
