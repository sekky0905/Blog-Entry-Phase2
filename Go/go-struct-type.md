# Goのstructの宣言にtypeをつける理由とGoのコードで見かけるstruct{}{}
2つの記事にしようかと思ったが、一緒の記事の方がわかりやすいと思ったので、2つのセクションに分けて1つの記事で記す。

## TL;DR
構造体にtypeをつけて名前をつけるのは、
* 色々な場面でその構造体を使用したいから
	* 構造体にtypeを付けないとその構造体を使用する度に、構造体を最初から宣言しなければならないから

struct{}{}は、
* 無名構造体である
* struct{}という型
* 後ろの方の{}は、フィールドを何も持たないという意味

## Goのstructの宣言にtypeをつける理由
普段、構造体にtypeを付与していると思うが、なぜそうするのかということをまとめてみた。

### よくある構造体の宣言
以下のようにstructにtypeで型名をつけて構造体を宣言することが多いと思う。

```go
type User struct {
	ID   string
	Name string
}
```

構造体にtypeを使う時は、あくまでも以下である。
> structで定義された構造体にtypeを使って新しい型名を与える

引用元: 松尾 愛賀　(2016/4/15)『スターティングGo言語』 翔泳社

typeとはそもそも以下のようなものだからだ。
> typeという予約語を用いると、既存の型や型リテラルに別名をつけることができます

引用元:[Goを学びたての人が誤解しがちなtypeと構造体について #golang - Qiita](https://qiita.com/tenntenn/items/45c568d43e950292bc31)

typeに関しての説明は、上記の記事に詳しく記述されているので、参照されたし。

### typeを使わないでフィールドを持った構造体を宣言してみる
typeはあくまでも、新しい型の名前をつけるだけなので、以下のようにフィールドを定義した無名構造体を使うようなことはできる。
この場合、uでのみ使える型になる。

```go
u := struct {
	ID string
	Name string
}{ID: "hogehoge", Name: "太郎"}
fmt.Printf("uの型は%T\n", u)
```

実行結果

```
uの型はstruct { ID string; Name string }
```

この場合、同じ型を別の変数に格納したい場合は、もう一回以下の部分を書かなきゃいけなくて冗長である。

```go
struct {
	ID   string
	Name string
}
```

type MyInt intのように、既存の型に別名をつけるのと何ら変わらないが、intとかの基本型よりも構造体の方がtypeをつけるのが多いのは、構造体にはフィールドが存在するので、当たり前だけど、そのフィールドを持った構造体をいちいち宣言するのは冗長だから。
構造体にtypeをつけて名前をつけるのは、色々な場面でその構造体を使用したいからである。

## Goのコードで見かけるstruct{}{}
上記で、typeを使わないでフィールドを持った構造体を使用する例を見てきた。
Goのコードで、struct{}{}というコードを見かけたことはないだろうか。
以降は、このstruct{}{}について記述する。

## struct{}{}
TL;DRに書いた通り、これは、無名のstruct{}という型で、かつ後ろの方の{}は、フィールドを何も持たないよという意味。

コードで見てみると以下のようになる。

```go
hoge := struct{}{}
fmt.Printf("hogeの型は%T\n", hoge)
```

実行結果

```
hogeの型はstruct {}
```

ちなみに、struct{}にtypeで名前を付けると以下のようになる。

```go
func main() {
	type User struct {
	}

	hoge := &User{}
	fmt.Printf("hogeの型は%T\n", hoge)
}
```

実行結果

```
hogeの型は*main.User
```

以下のような場合は、当たり前だがエラーが出る。

```go
u := struct{}{ID: "hogehoge", Name: "太郎",}
```

なぜなら、以下のエラー内容にあるようにstruct {}の型にはそのようなフィールドが存在しないからだ。

エラー

```
unknown field 'ID' in struct literal of type struct {}
unknown field 'Name' in struct literal of type struct {}
```

## 参考書籍
松尾 愛賀　(2016/4/15)『スターティングGo言語』 翔泳社

## 参考にさせていただいたサイト
[Goを学びたての人が誤解しがちなtypeと構造体について #golang - Qiita](https://qiita.com/tenntenn/items/45c568d43e950292bc31)
