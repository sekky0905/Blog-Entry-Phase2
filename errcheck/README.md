# Goでエラーハンドリングが行われているかを静的に解析する
Goでエラーハンドリングをしているかどうかを静的に解析してくれる[kisielk/errcheck](https://github.com/kisielk/errcheck)というツールが便利なので、使い方を簡単にまとめた。

[kisielk/errcheck: errcheck checks that you checked errors.](https://github.com/kisielk/errcheck/blob/master/README.md)を大変参考にさせていただいた。

# 基本的な機能
基本的な機能としては、エラーをチェックしているか(ハンドリング)しているかどうかをチェックしてくれる。
例えば以下のようなコードがあった場合に、mainでvalidate関数を呼び出しているが、エラーハンドリングを行っていないため、`errcheck ./...` で下記のようなメッセージが表示される。

コード

```go
package main

import (
	"errors"
	"fmt"
)

func main() {
	// ここでエラーハンドリングを本来行うべきだが、行われていないのでメッセージが表示される。
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

以下のようにすれば、上記のメッセージは表示されなくなる。

```go
if err := validate(19); err != nil {
	// エラーハンドリング
}
```

# インストール
`go get -u github.com/kisielk/errcheck`

# エラーチェックを行う対象の指定の仕方
以下のコマンドは、
[公式のREADME.md](https://github.com/kisielk/errcheck/blob/master/README.md)から引用。

## 特定のpackageをチェックする
`package errcheck github.com/path/to/package`

## カレントディレクトリ配下全部をチェックする
`errcheck ./...`

##  $GOPATH and $GOROOT配下の全部のpackageをチェックする
`errcheck all`

# その他の機能
## エラーチェックしたくないものが存在する場合

エラーチェックしたくないものもあるかもしれない。
そういう場合には、エラーチェックをしたくない関数のリストを記述したファイル(以下のコマンドでは`errcheck_excludes.txt`)を作成し、以下のようにコマンドを叩く。

> errcheck -exclude errcheck_excludes.txt path/to/package

引用元 : https://github.com/kisielk/errcheck 

また、
`errcheck -ignore 'Close' ./...` のように引数を与えると、Closeの部分のみ無視するようにすることもできる。

# オプション
kisielk/errcheckには、いくつかオプションが存在する。

## -tag
tagsオプションを用いてエラーチェックを行っているか確認する対象のファイルを `build tags` によって切り替えることができる。

例えば、以下のようなコードが存在するとする。

コード

```go:main.go
package main

func main() {
	caller()
}
```

```go:sample1.go
// +build sample1

package main

import (
	"errors"
	"fmt"
)

func caller(){
	worker()
}

func worker() error {
	fmt.Println("sample1")
	return errors.New("エラーだよ")
}
```

```go:sample2.go
// +build !sample1

package main

import "errors"

func caller(){
	worker()
}

func worker() error {
	fmt.Println("sample2")
	return errors.New("エラーだよ")
}
```

Goでは `build tags` を使う上記のようなコードの場合、`go build -tags sample1` のような感じでビルドすると、 `+build !sample1` が記述されたsample1.goがビルドされる。

参考 : [goで#ifdefのような条件分岐コンパイル - Qiita](https://qiita.com/yamasaki-masahide/items/8e5d59467dcf67d9b0be)
[Goで任意のbuild tagsをつけてビルドファイルを切り替える - Qiita](https://qiita.com/ueokande/items/fac0d1219dbbc8f44db7)

errcheckでも、tagsオプションを用いてエラーチェックを行っているか確認する対象のファイルを `build tags` で切り替えることができる。

上記のようなコードで `errcheck -tags sample1 ./...` とすると以下のようなメッセージが表示される。

メッセージ

```
sample1.go:11:8:	worker()
```

また、普通に `errcheck  ./...` とすると以下のようなメッセージが表示される。

```
sample2.go:11:8:	worker()
```

## -asserts
Goでは、 `t, ok := i.(T)` のような感じで型アサーションを行い、第二引数のokの部分で型アサーションが成功したかどうかを確認するとこができる。
このオプションでは、そのokでの確認が行われているかどうかをチェックしてくれる。
ちなみにこのokの部分をチェックしていないで、型アサーションが失敗するとpanicを起こす。

例えば以下のようなコードがあった場合に、Bar関数で型アサーションを使用しているが、第二引数のokの部分で型アサーションが成功したかどうかの確認を行っていないため、 `errcheck -asserts ./...` で下記のようなメッセージが表示される。

コード

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
	// 型アサーションが成功したかどうかのチェックを行っていない
	hoge := arg.(Hoge)
	hoge.Method("test")
}
```

メッセージ

```
main.go:27:10:	hoge := arg.(Hoge)
```

以下のように、二番目の戻り値で型変換が可能かどうかのチェックを行う必要があるので以下のようにするとメッセージは表示されなくなる。

```go
func Bar(arg interface{}) {
	hoge, ok := arg.(Hoge)
	if ok {
		// 何かする
	} else {
		// 何かする
	}
	hoge.Method("test")
}
```

## -blank
Goでは、関数やメソッドの戻り値を `_` で受け取って無視することができる。
例えば、`_ = method(arg) ` のようにだ。
エラーを返す関数やメソッドにおいて単純に `errcheck ./...` としても上記のような場合は、メッセージが表示されることはない。
しかし、これではエラーハンドリングを行っていないも同然になってしまう。
その場合、 、`errcheck -blank ./...` と `-blank` を付与すると上記のようなメソッドや関数が返すerrorを `_` で受け取って無視する箇所を指摘するメッセージを表示してくれる。

例えば以下のようなコードがあった場合に、mainでvalidate関数を呼び出しているが、validateが返すerrorを`_` で受け取って無視しているため、`errcheck -blank ./...` で下記のようなメッセージが表示される。

コード

```go
package main

import (
"errors"
"fmt"
)

func main() {
	// エラーを '_' で受け取って無視している。
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

メッセージ

```
blank_sample/main.go:10:2:	_ = validate(20)
```

以下のようにすれば、上記のメッセージは表示されなくなる。

```go
if err := validate(20); err != nil {
	// エラーハンドリング
}
```

## -abspath

`errcheck -abspath ./...` のようにすると、以下のようにエラーチェックが行われていない箇所が存在するファイルの絶対パスも表示してくれる。

メッセージ

```
/absolete/path/main.go:7:9: hoge()
```


# 参考にさせていただいたURL
[errcheck/README.md at master · kisielk/errcheck](https://github.com/kisielk/errcheck/blob/master/README.md)

[goで#ifdefのような条件分岐コンパイル - Qiita](https://qiita.com/yamasaki-masahide/items/8e5d59467dcf67d9b0be)

[Goで任意のbuild tagsをつけてビルドファイルを切り替える - Qiita](https://qiita.com/ueokande/items/fac0d1219dbbc8f44db7)