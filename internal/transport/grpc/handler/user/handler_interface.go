package user

import (
	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"
)

type Handler = v1.UserServiceServer

var _ Handler = (*HandlerImpl)(nil)
