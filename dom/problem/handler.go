package problem

import "fmt"

type Key string

type Problem struct {
	Key     Key
	Message string
}

func (id Key) WithMessage(message string) *Problem {
	return &Problem{
		Key:     id,
		Message: message,
	}
}

func (p *Problem) Error() string {
	return fmt.Sprintf("problem ocurred %s: %s", p.Key, p.Message)
}
