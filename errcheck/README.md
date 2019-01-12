# Goでエラーハンドリングが行われているかを静的に解析する

# 基本的な機能
基本的な機能としては、エラーをチェックしているか(ハンドリング)しているかどうかをチェックしてくれる。
例えば以下のようなコードがあった場合に、mainでvalidate関数を呼び出しているが、エラーハンドリングを行っていないため、下記のようなメッセージが表示される。

コード

```go
package main

import (
	"errors"
	"fmt"
)

func main() {
	// この状態でerrcheck main.goとすると
	// main.go:10:10:	validate(19)となる
	validate(19)
}

func validate(age int) error {
	if age < 20 {
		return errors.New("age should be 20 or more")
	}
	fmt.Println("ok~")
	return nil
}
```

メッセージ

```
main.go:10:10:	validate(19)
```

## エラーチェックしたくないものが存在する場合

エラーチェックしたくないものもあるかもしれない。
そういう場合には、エラーチェックをしたくない関数のリストを作成し、以下のように実装すると良い。

> errcheck -exclude errcheck_excludes.txt path/to/package

引用元 : https://github.com/kisielk/errcheck 

また、
`errcheck -ignore 'Close' ./...` のように引数を与えると、Closeの部分のみ無視するようにすることもできる。

# オプション
kisielk/errcheckには、いくつかオプションが存在する。

## -tag
custom build tagsの際に関係するtagを引数に与える。
`go build -tags hoge` のような時。

## -asserts
まずはコードを。

```go
package main

import "fmt"

type Hoge interface {
	Method(string)
}

type Foo struct {

}

func (f Foo)Method(arg string) {
	fmt.Println(arg)
}

func NewHoge()Hoge {
	return &Foo{}
}

func main() {
	hoge := NewHoge()
	Bar(hoge)
}

func Bar(arg interface{}) {
	hoge := arg.(Hoge)
	hoge.Method("test")
}
```

上記のようなコードの場合、本来はBarのところで、以下のように、二番目の戻り値で型変換が可能かどうかのチェックを行う必要がある。

```go
func Bar(arg interface{}) {
	hoge, ok := arg.(Hoge)
	if ok {
		// 何かする
	}
	hoge.Method("test")
}
```

実際にこのようなチェックを行っていないコードを `errcheck -asserts ./...` な形でチェックを行うと、`main.go:27:10:	hoge := arg.(Hoge)` のように知らせてくれるようになる。

## -blank

まずはコードを。

```go
package main

import (
"errors"
"fmt"
)

func main() {
	 _ = validate(20)
}

func validate(age int) error {
	if age < 20 {
		return errors.New("age should be 20 or more")
	}
	fmt.Println("ok~")
	return nil
}
```

このコードを単純に `errcheck ./...` のようにすると、チェック結果は特に問題を返さない。
`_` に `error` を格納し、何もしていないように見えるが、errorcheck的には、問題にしない。
このような場合に、`errcheck -blank ./...` を行うと、`main.go:9:3:	_ = validate(20)` のように問題の箇所を教えてくれる。

