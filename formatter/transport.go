package formatter

type ITransport interface {
	Write([]byte) error
	Close() error
}
