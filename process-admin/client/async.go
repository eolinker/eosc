package client

import (
	"context"
	"github.com/eolinker/eosc/process-admin/model"
)

type Response[T any] struct {
	Result T
	Err    error
}
type BoolResponse = Response[bool]
type ResponseChan[T any] chan<- *Response[T]

func MakeResponseChan[T any]() ResponseChan[T] {
	return make(chan *Response[T], 1)
}

type AsyncClient interface {
	List(ctx context.Context, profession string, c ResponseChan[[]model.Object]) error
	Get(ctx context.Context, id string, c ResponseChan[model.Object]) error
	Set(ctx context.Context, id string, value any, c ResponseChan[bool]) error
	Del(ctx context.Context, id string, c ResponseChan[bool]) error
	Exists(ctx context.Context, id string, c ResponseChan[bool]) (bool, error)
	PList(ctx context.Context, c ResponseChan[[]*model.ProfessionInfo]) error
	PGet(ctx context.Context, name string, c ResponseChan[*model.ProfessionConfig]) error
	SGet(ctx context.Context, name string, c ResponseChan[model.Object]) error
	SSet(ctx context.Context, name string, value any, c ResponseChan[bool]) error
	VAll(ctx context.Context, c ResponseChan[map[string]model.Variables]) error
	VGet(ctx context.Context, namespace string, c ResponseChan[model.Variables]) error
	VSet(ctx context.Context, namespace string, v model.Variables, c ResponseChan[bool]) error
	Begin(ctx context.Context, c ResponseChan[bool]) error
	Commit(ctx context.Context, c ResponseChan[bool]) error
	Rollback(ctx context.Context, c ResponseChan[bool]) error
}
