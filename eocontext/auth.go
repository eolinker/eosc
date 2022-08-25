package eocontext

type IAuthHandler interface {
	Auth(ctx EoContext) error
}
