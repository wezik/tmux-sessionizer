package utils

import "fmt"

// Well I like it
func Ensure(condition bool, message string) {
	if !condition {
		fmt.Println(message)
		panic(message)
	}
}

func EnsureWithErr(condition bool, err error) {
	if !condition {
		fmt.Println(err.Error())
		panic(err)
	}
}
