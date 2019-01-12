package main

import (
	"errors"
	"fmt"
)

func main() {
	// エラーハンドリングをblank
	_ = validate(20)
}

func validate(age int) error {
	if age < 20 {
		return errors.New("age should be 20 or more")
	}
	fmt.Println("ok~")
	return nil
}
