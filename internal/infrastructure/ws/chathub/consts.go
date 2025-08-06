package chathub

import "time"

const (
	writeWait      = 15 * time.Second
	pongWait       = 10 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
	chanBuff       = 1024
)
