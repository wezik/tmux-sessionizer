package errors

import (
	"fmt"
	"os"
)

func EnsureNotNil(err error, message string) {
	if err != nil {
		fmt.Println(message, ":", err.Error())
		os.Exit(1)
	}
}
