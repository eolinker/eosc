package extender

import "fmt"

const (
	StatusSuccess = iota
	StatusInit
	StatusDownloadFault
	StatusCheckFault
)

type Status struct {
	Group   string
	Project string
	Version string
	Status  int
}

func (s *Status) Name() string {
	return fmt.Sprint(s.Group, ":", s.Project)
}
