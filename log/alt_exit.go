package log

import (
	"fmt"
	"os"
)

var handlers = []func(){}

func runHandler(handler func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, "Error: Logrus exit handler error:", err)
		}
	}()

	handler()
}

func runHandlers() {
	for _, handler := range handlers {
		runHandler(handler)
	}
}

// Exit runs all the Logrus atexit handlers and then terminates the program using os.Exit(code)
func Exit(code int) {
	runHandlers()
	os.Exit(code)
}

// RegisterExitHandler appends a Logrus Exit handler to the list of handlers,
// call logrus.Exit to invoke all handlers. The handlers will also be invoked when
// any Fatal log entry is made.
//
// This method is useful when a caller wishes to use logrus to log a fatal
// message but also needs to gracefully shutdown. An example usecase could be
// closing database connections, or sending a alert that the application is
// closing.
func RegisterExitHandler(handler func()) {
	handlers = append(handlers, handler)
}

// DeferExitHandler prepends a Logrus Exit handler to the list of handlers,
// call logrus.Exit to invoke all handlers. The handlers will also be invoked when
// any Fatal log entry is made.
//
// This method is useful when a caller wishes to use logrus to log a fatal
// message but also needs to gracefully shutdown. An example usecase could be
// closing database connections, or sending a alert that the application is
// closing.
func DeferExitHandler(handler func()) {
	handlers = append([]func(){handler}, handlers...)
}
