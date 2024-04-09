package proto

// apinto resp protocol data type.
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

func ReplyTypeString(r ReplyType) string {
	switch r {
	case ErrorReply:
		return "ErrorReply"
	case StatusReply:
		return "StatusReply"
	case IntReply:
		return "IntReply"
	case StringReply:
		return "StringReply"
	case ArrayReply:
		return "ArrayReply"
	default:
		return "unknown"
	}
}
