package problem

type Key string

func (k Key) WithMsg(msg string) Problem {
	return Problem{Key: k, Message: msg}
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
