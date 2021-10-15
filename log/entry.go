package log

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"
)

var (

	// qualified package name, cached at first use
	thisPackageName string

	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth int

	// Used for caller information initialisation
	callerInitOnce sync.Once
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 2
)

func init() {

	// start at the bottom of the stack before the package-name cache is primed
	minimumCallerDepth = 1
}

// Defines the key when adding errors using WithError.
var ErrorKey = "error"

// An entry is the final or intermediate Logrus logging entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Trace, Debug,
// Info, Warn, Error, Fatal or Panic is called on it. These objects can be
// reused and passed around as much as you wish to avoid field duplication.
type Entry struct {

	// Contains all the fields set by the user.
	Data Fields

	// Time at which the log entry was created
	Time time.Time

	// Level the log entry was logged at: Trace, Debug, Info, Warn, Error, Fatal or Panic
	// This field will be set on entry firing and the value will be equal to the one in _Logger struct field.
	Level Level

	// Calling method, with package name
	Caller *runtime.Frame

	// Message passed to Trace, Debug, Info, Warn, Error, Fatal or Panic
	Message string
	Err     string
}

func (entry *Entry) HasCaller() (has bool) {
	return entry.Caller != nil
}

type EntryBuilder struct {
	logger *Logger
	// Contains all the fields set by the user.
	Data Fields
	// Time at which the log entry was created
	Time      time.Time
	prefix    string
	hasPrefix bool
	err       string
}

func (builder *EntryBuilder) Logln(level Level, args ...interface{}) {
	if builder.logger.IsLevelEnabled(level) {
		builder.log(level, sprintlnn(args...))
	}
}

func (builder *EntryBuilder) Log(level Level, args ...interface{}) {
	if builder.logger.IsLevelEnabled(level) {
		builder.log(level, fmt.Sprint(args...))
	}
}

func (builder *EntryBuilder) Logf(level Level, format string, args ...interface{}) {
	if builder.logger.IsLevelEnabled(level) {
		builder.log(level, fmt.Sprintf(format, args...))
	}
}

// Add an error as single field (using the key defined in ErrorKey) to the Entry.
func (builder *EntryBuilder) WithError(err error) Builder {
	return builder.WithField(ErrorKey, err)
}

// Add a single field to the Entry.
func (builder *EntryBuilder) WithField(key string, value interface{}) Builder {
	return builder.WithFields(Fields{key: value})
}

// Add a map of fields to the Entry.
func (builder *EntryBuilder) WithFields(fields Fields) Builder {
	data := make(Fields, len(builder.Data)+len(fields))
	for k, v := range builder.Data {
		data[k] = v
	}

	for k, v := range fields {
		isErrField := false
		if t := reflect.TypeOf(v); t != nil {
			switch t.Kind() {
			case reflect.Func:
				isErrField = true
			case reflect.Ptr:
				isErrField = t.Elem().Kind() == reflect.Func
			}
		}
		if !isErrField {
			data[k] = v
		}
	}
	return &EntryBuilder{logger: builder.logger, Time: builder.Time, Data: data}
}

// getCaller retrieves the name of the first non-logrus calling function
func getCaller(packageName string) *runtime.Frame {

	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, 2)
		_ = runtime.Callers(0, pcs)
		thisPackageName = getPackageName(runtime.FuncForPC(pcs[1]).Name())

		// now that we have the cache, we can skip a minimum count of known-logrus functions
		// XXX this is dubious, the number of frames may vary
		minimumCallerDepth = knownLogrusFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != thisPackageName && pkg != packageName {
			return &f
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// This function is not declared with a pointer value because otherwise
// race conditions will occur when using multiple goroutines
func (builder *EntryBuilder) log(level Level, msg string) error {
	defer builder.logger.pool.Put(builder)
	if builder.hasPrefix {
		msg = fmt.Sprint(builder.prefix, msg)
	}
	entry := &Entry{
		Data:    builder.Data,
		Time:    builder.Time,
		Level:   level,
		Caller:  nil,
		Message: msg,
		Err:     builder.err,
	}

	if builder.logger.reportCaller {
		entry.Caller = getCaller(builder.logger.packageName)
	}
	return builder.logger.Transport(entry)

}

// Sprintlnn => Sprint no newline. This is to get the behavior of how
// fmt.Sprintln where spaces are always added between operands, regardless of
// their type. Instead of vendoring the Sprintln implementation to spare a
// string allocation, we do the simplest thing.
func sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}
