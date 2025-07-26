package vo

type AuthMiddlewareKey string

const (
	UserIDKey  AuthMiddlewareKey = "user_id"
	ChatIDsKey AuthMiddlewareKey = "chat_ids"
)
