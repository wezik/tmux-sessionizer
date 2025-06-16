package problem

import "fmt"

type Key string

func (k Key) WithMsg(a ...any) Problem {
	return Problem{Key: k, Message: fmt.Sprint(a...)}
}

type Problem struct {
	Key     Key
	Message string
}

func (p Problem) Error() string {
	return p.Message
}

func (p Problem) EqualKey(other error) bool {
	return p.Key == other.(Problem).Key
}
