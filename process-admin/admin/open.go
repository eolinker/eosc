package admin

import "time"

func GenVersion() string {
	return time.Now().Format("20060102150405")
}
