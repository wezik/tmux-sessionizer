package errors

import (
	"fmt"
	"os"
)

func EnsureNotNil(err error, message string) {
	if err != nil {
		fmt.Println(message)
		os.Exit(1)
	}
}
