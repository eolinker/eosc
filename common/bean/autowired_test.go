package bean

import (
	"testing"
)

type Tester interface {
	GET() string
}
type Tester1 struct {
}
type Tester2 struct {
}

func (t *Tester2) GET() string {
	return "test2"
}

func (t *Tester1) GET() string {
	return "test1"
}

func Test(t *testing.T) {

	factory := NewContainer()
	var t1 Tester = new(Tester1)

	factory.Injection(&t1)

	var tester Tester = nil
	factory.Autowired(&tester)

	if tester != nil {
		r := tester.GET()
		if r != t1.GET() {
			t.Error("unknown error")
		}
	}

	if err := factory.Check(); err != nil {
		t.Error(err)
		return
	}
	r := tester.GET()
	if r != t1.GET() {
		t.Error("unknown error")
	}

}
