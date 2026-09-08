package chat

import (
	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"
)

type Handler = v1.ChatServiceServer

var _ Handler = (*HandlerImpl)(nil)
