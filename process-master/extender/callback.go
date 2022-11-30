package extender

type ICallback interface {
	Update(es []*Status, success bool)
}

type CallbackList []ICallback

func GenCallbackList(hs ...ICallback) ICallback {
	return CallbackList(hs)
}
func (ecs CallbackList) Update(es []*Status, success bool) {
	for _, e := range ecs {
		e.Update(es, success)
	}
}
