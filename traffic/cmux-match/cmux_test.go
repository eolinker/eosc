// Copyright 2016 The CMux Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package cmuxMatch

import (
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/soheilhy/cmux"
	"go/build"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	testHTTP1Resp = "http1"
	rpcVal        = 1234
)

func safeDial(t *testing.T, addr net.Addr) (*rpc.Client, func()) {
	c, err := rpc.Dial(addr.Network(), addr.String())
	if err != nil {
		t.Fatal(err)
	}
	return c, func() {
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

type chanListener struct {
	net.Listener
	connCh chan net.Conn
}

func newChanListener() *chanListener {
	return &chanListener{connCh: make(chan net.Conn, 1)}
}

func (l *chanListener) Accept() (net.Conn, error) {
	if c, ok := <-l.connCh; ok {
		return c, nil
	}
	return nil, errors.New("use of closed network connection")
}

func testListener(t *testing.T) (net.Listener, func()) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	var once sync.Once
	return l, func() {
		once.Do(func() {
			if err := l.Close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

type testHTTP1Handler struct{}

func (h *testHTTP1Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, testHTTP1Resp)
}

func runTestHTTPServer(errCh chan<- error, l net.Listener) {
	var mu sync.Mutex
	conns := make(map[net.Conn]struct{})

	defer func() {
		mu.Lock()
		for c := range conns {
			if err := c.Close(); err != nil {
				errCh <- err
			}
		}
		mu.Unlock()
	}()

	s := &http.Server{
		Handler: &testHTTP1Handler{},
		ConnState: func(c net.Conn, state http.ConnState) {
			mu.Lock()
			switch state {
			case http.StateNew:
				conns[c] = struct{}{}
			case http.StateClosed:
				delete(conns, c)
			}
			mu.Unlock()
		},
	}
	if err := s.Serve(l); err != cmux.ErrListenerClosed && err != cmux.ErrServerClosed {
		errCh <- err
	}
}

func generateTLSCert(t *testing.T) {
	err := exec.Command("go", "run", build.Default.GOROOT+"/src/crypto/tls/generate_cert.go", "--host", "*").Run()
	if err != nil {
		t.Fatal(err)
	}
}

func cleanupTLSCert(t *testing.T) {
	err := os.Remove("cert.pem")
	if err != nil {
		t.Error(err)
	}
	err = os.Remove("key.pem")
	if err != nil {
		t.Error(err)
	}
}

func runTestTLSServer(errCh chan<- error, l net.Listener) {
	certificate, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		errCh <- err
		log.Printf("1")
		return
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		Rand:         rand.Reader,
	}

	tlsl := tls.NewListener(l, config)
	runTestHTTPServer(errCh, tlsl)
}

func runTestHTTP1Client(t *testing.T, addr string) {
	runTestHTTPClient(t, "http", addr)
}

func runTestTLSClient(t *testing.T, addr string) {
	runTestHTTPClient(t, "https", addr)
}

func runTestHTTPClient(t *testing.T, proto string, addr string) {
	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	r, err := client.Get(proto + "://" + addr)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = r.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != testHTTP1Resp {
		t.Fatalf("invalid response: want=%s got=%s", testHTTP1Resp, b)
	}
}

type TestRPCRcvr struct{}

func (r TestRPCRcvr) Test(i int, j *int) error {
	*j = i
	return nil
}

func runTestRPCServer(errCh chan<- error, l net.Listener) {
	s := rpc.NewServer()
	if err := s.Register(TestRPCRcvr{}); err != nil {
		errCh <- err
	}
	for {
		c, err := l.Accept()
		if err != nil {
			if err != cmux.ErrListenerClosed && err != cmux.ErrServerClosed {
				errCh <- err
			}
			return
		}
		go s.ServeConn(c)
	}
}

func runTestRPCClient(t *testing.T, addr net.Addr) {
	c, cleanup := safeDial(t, addr)
	defer cleanup()

	var num int
	if err := c.Call("TestRPCRcvr.Test", rpcVal, &num); err != nil {
		t.Fatal(err)
	}

	if num != rpcVal {
		t.Errorf("wrong rpc response: want=%d got=%v", rpcVal, num)
	}
}

const (
	handleHTTP1Close   = 1
	handleHTTP1Request = 2
	handleAnyClose     = 3
	handleAnyRequest   = 4
)

func TestTimeout(t *testing.T) {
	defer leakCheck(t)()
	lis, Close := testListener(t)
	defer Close()
	result := make(chan int, 5)
	testDuration := time.Millisecond * 500
	m := NewMatch(lis)
	m.SetReadTimeout(testDuration)
	http1 := m.Match(Http1)
	any := m.Match(Any)

	go func() {
		con, err := http1.Accept()
		if err != nil {
			result <- handleHTTP1Close
		} else {
			_, _ = con.Write([]byte("http1"))
			_ = con.Close()
			result <- handleHTTP1Request
		}
	}()
	go func() {
		con, err := any.Accept()
		if err != nil {

			result <- handleAnyClose
		} else {
			_, _ = con.Write([]byte("any"))
			_ = con.Close()
			result <- handleAnyRequest
		}
	}()
	time.Sleep(testDuration) // wait to prevent timeouts on slow test-runners
	client, err := net.Dial("tcp", lis.Addr().String())
	if err != nil {
		log.Fatal("testTimeout client failed: ", err)
	}
	defer func() {
		_ = client.Close()
	}()
	time.Sleep(testDuration / 2)
	if len(result) != 0 {
		log.Print("tcp ")
		t.Fatal("testTimeout failed: accepted to fast: ", len(result))
	}
	_ = client.SetReadDeadline(time.Now().Add(testDuration * 3))
	buffer := make([]byte, 10)
	rl, err := client.Read(buffer)
	if err != nil {
		t.Fatal("testTimeout failed: client error: ", err, rl)
	}
	Close()
	if rl != 3 {
		log.Print("testTimeout failed: response from wrong sevice ", rl)
	}
	if string(buffer[0:3]) != "any" {
		log.Print("testTimeout failed: response from wrong sevice ")
	}
	time.Sleep(testDuration * 2)
	if len(result) != 2 {
		t.Fatal("testTimeout failed: accepted to less: ", len(result))
	}
	if a := <-result; a != handleAnyRequest {
		t.Fatal("testTimeout failed: any rule did not match")
	}
	if a := <-result; a != handleHTTP1Close {
		t.Fatal("testTimeout failed: no close an http rule")
	}
}

func TestAny(t *testing.T) {
	defer leakCheck(t)()
	errCh := make(chan error)
	defer func() {
		select {
		case err := <-errCh:
			t.Fatal(err)
		default:
		}
	}()
	l, cleanup := testListener(t)
	defer cleanup()

	muxl := NewMatch(l)
	httpl := muxl.Match(Any)

	go runTestHTTPServer(errCh, httpl)

	runTestHTTP1Client(t, l.Addr().String())
}

func TestTLS(t *testing.T) {
	generateTLSCert(t)
	defer cleanupTLSCert(t)
	defer leakCheck(t)()
	errCh := make(chan error)
	defer func() {
		select {
		case err := <-errCh:
			t.Fatal(err)
		default:
		}
	}()
	l, cleanup := testListener(t)
	defer cleanup()

	muxl := NewMatch(l)
	tlsl := muxl.Match(Https)
	httpl := muxl.Match(Any)

	go runTestTLSServer(errCh, tlsl)
	go runTestHTTPServer(errCh, httpl)

	runTestHTTP1Client(t, l.Addr().String())
	runTestTLSClient(t, l.Addr().String())
}

func TestHTTPGoRPC(t *testing.T) {
	defer leakCheck(t)()
	errCh := make(chan error)
	defer func() {
		select {
		case err := <-errCh:
			t.Fatal(err)
		default:
		}
	}()
	l, cleanup := testListener(t)
	defer cleanup()

	muxl := NewMatch(l)
	http2l := muxl.Match(Http2)
	httpl := muxl.Match(Http1)
	rpcl := muxl.Match(Any)

	go runTestHTTPServer(errCh, http2l)
	go runTestHTTPServer(errCh, httpl)
	go runTestRPCServer(errCh, rpcl)

	runTestHTTP1Client(t, l.Addr().String())
	runTestRPCClient(t, l.Addr())
}

// Cribbed from google.golang.org/grpc/test/end2end_test.go.

// interestingGoroutines returns all goroutines we care about for the purpose
// of leak checking. It excludes testing or runtime ones.
func interestingGoroutines() (gs []string) {
	buf := make([]byte, 2<<20)
	buf = buf[:runtime.Stack(buf, true)]
	for _, g := range strings.Split(string(buf), "\n\n") {
		sl := strings.SplitN(g, "\n", 2)
		if len(sl) != 2 {
			continue
		}
		stack := strings.TrimSpace(sl[1])
		if strings.HasPrefix(stack, "testing.RunTests") {
			continue
		}

		if stack == "" ||
			strings.Contains(stack, "main.main()") ||
			strings.Contains(stack, "testing.Main(") ||
			strings.Contains(stack, "runtime.goexit") ||
			strings.Contains(stack, "created by runtime.gc") ||
			strings.Contains(stack, "interestingGoroutines") ||
			strings.Contains(stack, "runtime.MHeap_Scavenger") {
			continue
		}
		gs = append(gs, g)
	}
	sort.Strings(gs)
	return
}

// leakCheck snapshots the currently-running goroutines and returns a
// function to be run at the end of tests to see whether any
// goroutines leaked.
func leakCheck(t testing.TB) func() {
	orig := map[string]bool{}
	for _, g := range interestingGoroutines() {
		orig[g] = true
	}
	return func() {
		// Loop, waiting for goroutines to shut down.
		// Wait up to 5 seconds, but finish as quickly as possible.
		deadline := time.Now().Add(5 * time.Second)
		for {
			var leaked []string
			for _, g := range interestingGoroutines() {
				if !orig[g] {
					leaked = append(leaked, g)
				}
			}
			if len(leaked) == 0 {
				return
			}
			if time.Now().Before(deadline) {
				time.Sleep(50 * time.Millisecond)
				continue
			}
			for _, g := range leaked {
				t.Errorf("Leaked goroutine: %v", g)
			}
			return
		}
	}
}
