package chat

import "time"

const (
	maxMessageSize = 512 * 1024
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	writeWait      = 10 * time.Second
)
