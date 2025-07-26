package vo

type OutboxStatus string

const (
	StatusPending   OutboxStatus = "pending"
	StatusProcessed OutboxStatus = "processed"
	StatusFailed    OutboxStatus = "failed"
)
