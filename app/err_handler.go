package main

import (
	"fmt"
	"os"
	"thop/dom/problem"
	"thop/dom/selector"
)

func handleErr(err error) {
	if err == nil {
		return
	}

	pb := err.(*problem.Problem)
	if pb == nil {
		panic(err)
	}

	if pb.Key == selector.ErrCancelled {
		fmt.Println(pb.Message)
		return
	}

	fmt.Println(err)
	os.Exit(1)
}
