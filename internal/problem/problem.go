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

// Compare only by keys, not by message
func (k Key) Equal(other error) bool {
	if other == nil {
		return false
	}
	return k == other.(Problem).Key
}
