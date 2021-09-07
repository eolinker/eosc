package utils

import (
	"fmt"
	"net"
	"regexp"
)

func ValidAddr(addr string) bool {
	match, err := regexp.MatchString(`^(?:(?:1[0-9][0-9]\.)|(?:2[0-4][0-9]\.)|(?:25[0-5]\.)|(?:[1-9][0-9]\.)|(?:[0-9]\.)){3}(?:(?:1[0-9][0-9])|(?:2[0-4][0-9])|(?:25[0-5])|(?:[1-9][0-9])|(?:[0-9]))\:(([0-9])|([1-9][0-9]{1,3})|([1-6][0-9]{0,4}))$`, addr)
	if err != nil {
		return false
	}
	return match
}

func IsListen(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("the address %s is listened", addr)
	}
	listener.Close()
	return nil
}
