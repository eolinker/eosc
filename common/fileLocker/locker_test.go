package fileLocker

import "testing"

func TestLocker_Lock(t *testing.T) {
	masterLocker1 := NewLocker("", 10, "master")
	//masterLocker2 := NewLocker("", 10, "master")

	masterLocker1.Lock()
	//masterLocker2.Lock()
	cliLocker1 := NewLocker("", 10, "cli")
	//cliLocker2 := NewLocker("", 10, "cli")
	cliLocker1.Lock()
	//cliLocker2.Lock()
}
