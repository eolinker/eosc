package raft_service

//var ff = eosc.NewUntyped()

type ICreateHandler interface {
	Namespace() string
	Handler() IRaftServiceHandler
}

type createHandler struct {
	namespace string
	handler   IRaftServiceHandler
}

func (c *createHandler) Namespace() string {
	return c.namespace
}

func (c *createHandler) Handler() IRaftServiceHandler {
	return c.handler
}

func NewCreateHandler(namespace string, handler IRaftServiceHandler) ICreateHandler {
	return &createHandler{namespace: namespace, handler: handler}
}
