package utils

import (
	"bytes"
	"github.com/soheilhy/cmux"
	"io"
	"net"
)

var defaultHTTPMethods = []string{
	"OPTIONS",
	"GET",
	"HEAD",
	"POST",
	"PUT",
	"DELETE",
	"TRACE",
	"CONNECT",
}
var (
	methods         map[string]struct{}
	maxMethodLength int
)

func init() {
	methods = make(map[string]struct{})
	for _, m := range defaultHTTPMethods {
		b := []byte(m)
		if len(b) > maxMethodLength {
			maxMethodLength = len(b)
		}
		methods[m] = struct{}{}
	}
	maxMethodLength += 1
}
func MatchMux(ln net.Listener, paths ...[]string) []net.Listener {
	if len(paths) == 0 {
		return []net.Listener{ln}
	}

	mux := cmux.New(ln)
	rs := make([]net.Listener, 0, len(paths)+1)
	for _, path := range paths {
		rs = append(rs, mux.Match(HttpPrefixMatcher(path...)))
	}
	rs = append(rs, mux.Match(cmux.Any()))
	go func() {
		mux.Serve()
	}()
	return rs
}
func MatchForHttp(ln net.Listener, path ...string) (net.Listener, net.Listener) {
	mux := cmux.New(ln)
	http := mux.Match(HttpPrefixMatcher(path...))
	any := mux.Match(cmux.Any())

	go func() {
		mux.Serve()
	}()

	return http, any

}

func HttpPrefixMatcher(paths ...string) cmux.Matcher {
	maxPathLength := 0
	pathMaps := map[string]struct{}{}
	for _, p := range paths {
		if len(p) > maxPathLength {
			maxPathLength = len(p)
		}
		pathMaps[p] = struct{}{}
	}
	return func(reader io.Reader) bool {
		buf := make([]byte, maxMethodLength)
		n, err := io.ReadFull(reader, buf)
		if err != nil {
			return false
		}
		buf = buf[:n]
		indexByte := bytes.IndexByte(buf, ' ')
		if indexByte <= 0 {
			return false
		}
		if _, ok := methods[string(buf[:indexByte])]; !ok {
			return false
		}

		pb := make([]byte, maxPathLength)
		tem := buf[indexByte+1:]

		copy(pb, tem)
		np, _ := io.ReadFull(reader, pb[len(tem):])
		pb = pb[:np+len(tem)]
		_, has := pathMaps[string(pb)]
		return has

	}
}
