package responses

import "encoding/json"

type GetKeyResponse struct {
    Key              string `json:"key"`
    ValidationWindow string `json:"validation_window"`
}

type ValidateResponse struct {
    Status     string          `json:"status"`
    DataType   *string         `json:"data_type,omitempty"`
    DataValue  *string         `json:"data_value,omitempty"`
    PhoneData  json.RawMessage `json:"phone_data,omitempty"`
    DeviceData json.RawMessage `json:"device_data,omitempty"`
}
