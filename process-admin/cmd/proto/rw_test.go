package proto

import (
	"bufio"
	"io"
	"os"
	"testing"
	"time"
)

func TestRW(t *testing.T) {

	r, w, err := os.Pipe()
	if err != nil {
		return
	}

	go func() {
		wt := bufio.NewWriter(w)

		testWriter(wt, t)
		wt.Flush()
		w.Close()
	}()
	testReader(r, t)

}
func testReader(rt io.Reader, t *testing.T) {
	reader := NewReader(rt)

	for {
		message, err := reader.ReadMessage()
		if err != nil {
			if err == io.EOF {
				return
			}
			t.Error(err)
			return
		}

		printMessage(message, t)

	}

}
func printMessage(m IMessage, t *testing.T) {
	switch m.Type() {
	case ArrayReply:
		arr, err := m.Array()
		if err != nil {
			t.Error(err)
		}
		for _, i := range arr {
			printMessage(i, t)
		}
	default:
		s, err := m.String()
		if err != nil {
			return
		}
		t.Logf("read:%s", s)
	}
}
func testWriter(wt writer, t *testing.T) {
	w := NewWriter(wt)
	cmd1 := []any{
		"SetProvider", "api@router", map[string]any{
			"name":        "api@router",
			"type":        "router",
			"version":     "1.0.0",
			"description": "api@router",
		},
	}
	err := w.WriteArgs(cmd1...)
	if err != nil {
		t.Error(err)
		return
	}
	cmd2 := []any{
		"GET", "api@router",
	}
	err = w.WriteArgs(cmd2...)
	if err != nil {
		t.Error(err)
		return
	}
	cmd3 := []any{
		"response\r\ntest", time.Now(),
	}
	err = w.WriteArgs(cmd3...)
	if err != nil {
		t.Error(err)
		return
	}

}
