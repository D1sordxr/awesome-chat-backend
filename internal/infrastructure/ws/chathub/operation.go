package chathub

import (
	"context"
	"encoding/json"
)

type Operation struct {
	ClientID string
	Data     []byte
	Retries  int

	RespChan chan<- OperationResponse
	Ctx      context.Context
}

type OperationResponse struct {
	ID            int    `json:"id,omitempty"`
	OperationType string `json:"operation_type"`
	Success       bool   `json:"success"`
	Data          any    `json:"data,omitempty"` // []byte
	Error         error  `json:"error,omitempty"`
}

func (o *OperationResponse) ToJSON() []byte {
	payload, _ := json.Marshal(o)
	return payload
}

func SuccessResponse(opType string, data any) OperationResponse {
	//jsonData, _ := json.Marshal(data)
	return OperationResponse{
		OperationType: opType,
		Success:       true,
		Data:          data,
	}
}

func ErrorResponse(opType string, err error) OperationResponse {
	return OperationResponse{
		OperationType: opType,
		Success:       false,
		Error:         err,
	}
}

type operationHandler interface {
	Handle(ctx context.Context, data []byte) OperationResponse
}
