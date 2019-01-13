// +build !sample1

package main

import (
	"errors"
	"fmt"
)

func caller() {
	worker()
}

func worker() error {
	fmt.Println("sample2")
	return errors.New("エラーだよ")
}
