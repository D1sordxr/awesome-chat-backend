package user

//go:generate gowrap gen -g -p . -i Handler -t $GOWRAP_OP_TPL -o handler_with_tracing.go -v "OpPrefix=user.Handler"

type HandlerImpl struct {
	uc useCase
}

func NewHandler(uc useCase) *HandlerImpl {
	return &HandlerImpl{uc: uc}
}
