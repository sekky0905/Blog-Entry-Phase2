package main

import "fmt"

type Hoge interface {
	Method(string)
}

type Foo struct {
}

func (f Foo) Method(arg string) {
	fmt.Println(arg)
}

func NewHoge() Hoge {
	return &Foo{}
}

func main() {
	hoge := NewHoge()
	Bar(hoge)
}

func Bar(arg interface{}) {
	// hoge, ok := arg.(Hoge)で、okを確認した方が安全
	hoge := arg.(Hoge)
	hoge.Method("test")
}
