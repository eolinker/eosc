package proto

// redis resp protocol data type.
type ReplyType = byte

const (
	ErrorReply  = '-' // error
	StatusReply = '+' // string
	IntReply    = ':' // int
	StringReply = '$' // text
	ArrayReply  = '*' //array
)
const (
	Nil = Error("nil")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
