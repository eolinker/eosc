package log

// The Formatter interface is used to implement a custom Formatter. It takes an
// `Entry`. It exposes all the fields, including the default ones:
//
// * `entry.data["msg"]`. The message passed from DriverInfo, Warn, Error ..
// * `entry.data["time"]`. The timestamp.
// * `entry.data["level"]. The level the entry was logged at.
//
// Any additional fields added with `WithField` or `WithFields` are also in
// `entry.data`. Format is expected to return an array of bytes which are then
// logged to `logger.Out`.
type Formatter interface {
	Format(*Entry) ([]byte, error)
}
