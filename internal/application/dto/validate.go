package dto

import "encoding/json"

type ValidateInput struct {
	Nonce string
}

type ValidateOutput struct {
	Status     string
	DataType   *string
	DataValue  *string
	PhoneData  json.RawMessage
	DeviceData json.RawMessage
}
