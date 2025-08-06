package transport

import "encoding/json"

type OperationDTO struct {
	ID        int             `json:"id,omitempty"`
	Operation string          `json:"operation"`
	Body      json.RawMessage `json:"body"`
}
