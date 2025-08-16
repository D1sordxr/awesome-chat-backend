package transport

import "encoding/json"

type OperationHeader struct {
	ID        int             `json:"id,omitempty"`
	Operation string          `json:"operation"`
	Body      json.RawMessage `json:"body"`
}
