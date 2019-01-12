package main

import (
	"errors"
	"fmt"
)

func main() {
	// この状態でerrcheck main.goとすると
	// main.go:10:10:	validate(19)となる
	// validate(19)
	if err := validate(20); err != nil {
		fmt.Println(err)
	}
}

func validate(age int) error {
	if age < 20 {
		return errors.New("age should be 20 or more")
	}
	fmt.Println("ok~")
	return nil
}

