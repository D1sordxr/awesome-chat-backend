package message

import (
	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"
)

type Handler = v1.MessageServiceServer

var _ Handler = (*HandlerImpl)(nil)
