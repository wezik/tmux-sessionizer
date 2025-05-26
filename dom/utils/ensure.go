package utils

// Well I like it
func Ensure(condition bool, message string) {
	if !condition {
		panic(message)
	}
}

func EnsureWithErr(condition bool, err error) {
	if !condition {
		panic(err)
	}
}
