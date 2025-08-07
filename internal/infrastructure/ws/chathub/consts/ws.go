package consts

import "time"

const (
	WriteWait      = 15 * time.Second
	PongWait       = 10 * time.Second
	PingPeriod     = (PongWait * 9) / 10
	MaxMessageSize = 512
	ChanBuff       = 1024
)
