package main

import (
	"fmt"
	"net"
)

func validAddr(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("the address %s is listened", addr)
	}
	listener.Close()
	return nil
}
