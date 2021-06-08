package log

import "fmt"

type ConsoleHook struct {
}

func (c *ConsoleHook) Write(p []byte) (n int, err error) {
	fmt.Println(string(p))
	return len(p), nil
}
