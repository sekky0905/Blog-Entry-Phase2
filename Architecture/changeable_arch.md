# はじめに
アーキテクチャや設計の書籍や記事、これまでの経験も踏まえ、学んだ事をここにまとめたい。(まだ、勉強中なので微妙なところもあるかもしれません。お気付きの点があればご指摘いただけるとありがたいです。)

参考文献や参考記事は、本当に良書、良記事で非常に参考にさせていただきました。

**生意気なタイトルにしてしまいましたが、自分への戒めということもあってこのタイトルにさせていただいたので、ご容赦ください。**

## ある共通した話題
設計やアーキテクチャについて書かれたの書籍や記事を読んでいく中で、言葉は違えどかなりの高確率で共通するテーマが存在した。
そう、それが **「変更に強くなろう」** といった趣旨のテーマだ。
アーキテクチャや設計に関する書籍や記事は様々な方法論で、これを実現しようとしていた。

## 今回のテーマと記事の構成
今回は、「変更に強くなろう」というテーマの中で、重要だと感じた概念や考え方をまとめて実際にそれがどう生かされているかということを考察していきたい。
「単体テスト」については、「変更に強い」コードを意識すると「単体テスト」がしやすいという副次的な効果も現れることがわかったので、サブテーマとして記述したい。
また、上記の2つを同時に満たしている思う最近何かと話題のクリーンアーキテクチャについても記述したい。

正直、「変更に強くなる」というのをテーマに書こうとしたのだが、テストもしやすくなるし、最近学んだクリーンアーキテクチャもそれらに関係があることなので、書きたいし...となって少し詰め込みすぎた感が否めないですが、ご容赦ください。

そのため、今回の記事はざっくり大きく以下の3つの編成にしている。<br>
・変更に強くなる編<br>
・単体テスト編<br>
・クリーンアーキテクチャ編

# 変更に強くなる編
ここでは変更に強くなるための概念等を紹介する。
(書籍等では、他にももっとたくさん紹介されているが、ここでは基本的な一部のみを紹介する)

## 共通性/可変性分析
これは、『オブジェクト指向のこころ (SOFTWARE PATTERNS SERIES)』という書籍で紹介されていた概念である。

簡単にまとめると、
**共通性分析**とは、問題の中の共通して変化しないところを見つけ出すこと。
**可変性分析**とは、問題の中の変化しやすい要素を見つけ出すこと。

以下の一文が非常にわかりやすい。
> これは問題領域のどこが流動的に要素になるのかを識別し(「共通性分析」)、そのあと、それらがどのように変化するのかを識別する(「可変性分析」)というものです。

引用元 : 
アラン・シャロウェイ (著), ジェームズ・R・トロット (著), 村上 雅章 (翻訳)　(2014/3/11)『オブジェクト指向のこころ (SOFTWARE PATTERNS SERIES)』 丸善出版 113ページ

さらに同書には、それを具体的にソフトウェアに落とし込んでいく方法が記述されている。

> 問題領域中の特定部分に流動的要素がある場合、共通性分析によってそれらをまとめる概念を定義できるわけです。
こういった概念は抽象クラスによって表現できます。そして可変性分析によって洗い出された流動的要素は、具象クラスになります。

引用元 : 
アラン・シャロウェイ (著), ジェームズ・R・トロット (著), 村上 雅章 (翻訳)　(2014/3/11)『オブジェクト指向のこころ (SOFTWARE PATTERNS SERIES)』 丸善出版 113ページ

自分なりに解釈すると、何かソフトウェアを設計する前には、そのソフトウェアによって解決する問題の中において、具体的な事象や物とその事象や物の抽象的なものを考え出す。
その具体的な事象や物は似たような部分がいくつかあって、それらに共通する振る舞いを集めた概念を見つけ出す。
そして、具体的な事象や物は具象クラスに、共通する振る舞いをインターフェースや抽象クラスに落とし込んで設計していくのが大事なのだと思った。変化する具体的な問題とその問題に共通する抽象的な問題に分けるのだ。
いわば、抽象と具体に分ける。

## 依存関係
コードには依存関係が存在する。
例えば、AがBを呼んでおり、BがCを呼んでいるといった場合、依存関係は、A=>B=>Cといった具合になる。
この場合、=>の向きは一方向である。しかし、場合によっては、A<=>B<=>Cといった具合に、矢印が双方向に向いている場合もある。これはAとBが互いに、BとCが互いに依存しあってしまっているのだ。これを循環依存という。
これはコードを複雑にしてしまう要因らしい。
Goだと、循環依存はエラーになるほどだ。

参考:
エリック・エヴァンス(著)、 今関 剛 (監修)、 和智 右桂  (翻訳) (2011/4/9)『エリック・エヴァンスのドメイン駆動設計 (IT Architects’Archive ソフトウェア開発の実践)』 翔泳社

Robert C.Martin (著)、 角 征典  (翻訳)、 高木 正弘 (翻訳)　(2018/7/27)『Clean Architecture 達人に学ぶソフトウェアの構造と設計』 KADOKAWA

## 結合度
そのモジュールが他のモジュールにどれほど依存しているか(他のモジュールからどれほど独立しているか)ということ。
モジュール毎の結合度が低ければ、低いほど他のモジュールが変更になっても、影響を受けないので良いと考えられてる。変更に強い設計にするためにはこれを意識する必要がある。結合度が低いことを疎結合と言ったりする。

参考 : [モジュール結合度とは - IT用語辞典 Weblio辞書](https://www.weblio.jp/content/%E3%83%A2%E3%82%B8%E3%83%A5%E3%83%BC%E3%83%AB%E7%B5%90%E5%90%88%E5%BA%A6)

## 一旦整理
コードには依存関係があることもわかった。依存関係は循環参照することなく、片方向の参照が好ましいという。では、Aの具象クラスがBの具象クラスに依存し、Bの具象クラスがCの具象クラスに依存するというのは、どうだろうか。

「共通性/可変性分析」のセクションで、具体的なこと(具象クラス)は、変化しやすいことを説明した。本記事の冒頭を思い出してほしい。数々の良書が「変化に強くなろう」と主張しているのにも関わらず、片方向とはいえ、変更されやすい具象クラスに依存するのは良いのだろうか。
A=>B=>Cという風に依存関係があった場合、どれも具象クラスなので、変化しやすい。例えば、Cに変化があったら、Bはどうなるだろうか。Bに変化があったら、Aはどうなるだろうか...
Bは、Cの変更に伴って、コードを変更しなくてはならないし、AもBの変更に伴ってコードを変更しなくてはならない...辛い...
そりゃあ、数々の良書が「変更に強くなろう」というわけだ。「変更に強くなろう」というのは、「ある変更に伴いドミノ倒しのように発生する数々のコードの変更に耐えられる精神的な強さを持とう!」と言っているのだろうか。いや違う。ある変更があっても、他の部分に変更を(極力)生じさせない方法論を提唱してくれている。
実際の方法論や考え方を見ていこう。

## インターフェースとポリモーフィズム
オブジェクト指向やデザインパターンを勉強していると必ず出てくするこの2つの言葉。
変更に強くなるためには、この2つ(言語によっては抽象クラスなども含む)をうまく使うことが大事なようだ。共通性/可変性分析の項目で、変化しやすい具体的な事象や物と、それらに共通的する変わらない抽象を見つけるという話をしたが、このインターフェースとポリモーフィズムというのは、それらをうまく扱うことを可能にしてくれる。

なお、この記事では、インターフェースとポリモーフィズム自体はある程度理解している前提で話を進めるので、それら自体の説明はあまりしない。
もしインターフェースやポリモーフィズムが怪しい場合は、以下の記事等を参照。

[オブジェクト指向と10年戦ってわかったこと - Qiita](https://qiita.com/tutinoco/items/6952b01e5fc38914ec4e#%E3%83%9D%E3%83%AA%E3%83%A2%E3%83%BC%E3%83%95%E3%82%A3%E3%82%BA%E3%83%A0)

[15分でわかる かんたんオブジェクト指向 - Qiita](https://qiita.com/koher/items/6878c80014992900add7#%E3%83%9D%E3%83%AA%E3%83%A2%E3%83%BC%E3%83%95%E3%82%A3%E3%82%BA%E3%83%A0)

インターフェースのポリモーフィズムの具体例に関しては、後ほど記述する。

### 共通性/可変性分析とインターフェースとポリモーフィズム
共通性分析において発見した共通的な振る舞いをまとめた抽象をインターフェースとして定義する。
可変性分析おいて発見した個別の具体的な物や事を具象クラスや構造体として定義する。

先ほどの「共通性/可変性分析」のセクションで、具体的なこと(具象クラス)は、変化しやすいが、抽象的なこと(インターフェースや抽象クラス)は変化しにくいことがわかった。「え、じゃあ、具象クラスに依存させるんじゃなくて、インターフェースに依存させちゃえば、いいんじゃね?」というのがこの考え方。誤解を恐れずにいえば、インターフェースは、具象クラスの共通の振る舞いを集めた抽象的なものだ。
具体的に使うのは、具象クラスだけれども、それを使用するクラスは、インターフェースに依存させちゃえば、依存は変わらない。実際にやることが変わるだけ。後ほど、実際にコードの例でみる。

## DIP(依存関係逆転の法則)

> ソフトウエアモジュールを疎結合に保つための特定の形式を指す用語. この原則に従うとソフトウェアの振る舞いを定義する上位レベルのモジュールから下位レベルモジュールへの従来の依存関係は逆転し、結果として下位レベルモジュールの実装の詳細から上位レベルモジュールを独立に保つことができるようになる. この原則で述べられていることは以下の２つである

> A. 上位レベルのモジュールは下位レベルのモジュールに依存すべきではない. 両方とも抽象（abstractions)に依存すべきである.

> B. 抽象は詳細に依存してはならない. 詳細が抽象に依存すべきである.

引用元 : [依存性逆転の原則 - Wikipedia](https://ja.wikipedia.org/wiki/%E4%BE%9D%E5%AD%98%E6%80%A7%E9%80%86%E8%BB%A2%E3%81%AE%E5%8E%9F%E5%89%87)

依存関係逆転の法則は、 `インターフェース` で `ポリモーフィズム` を用いて、モジュール間の結合度を緩やかにするためのもの。もっと具体的にいうと、別のレイヤーのクラスなどを使用するときには、その具象クラスを直接使うのではなく、そのインターフェースを参照しようねということ。AというクラスがBというクラスを利用するときに、Bを直接利用するのではなくて、Aの抽象(抽象クラスやインターフェイス)を利用するとBの実装に変更があっても左右されにくいので、そういう風にしましょうということ。

「共通性/可変性分析とインターフェースとポリモーフィズム」のセクションで記述した事を原則として切り出したものだ。
「変更に強くなる」とか、「単体テストをしやすくする」などの事を考えると、この原則は本当に重要なものだ。

なぜ依存関係の **逆転** というかは、この後の具体例のところでUMLぽいものを書いて説明する。

参考:
[クリーンアーキテクチャ(The Clean Architecture翻訳) | blog.tai2.net](https://blog.tai2.net/the_clean_architecture.html)

[オブジェクト指向設計原則とは - Qiita](https://qiita.com/UWControl/items/98671f53120ae47ff93a)

[依存性逆転の原則 - Wikipedia](https://ja.wikipedia.org/wiki/%E4%BE%9D%E5%AD%98%E6%80%A7%E9%80%86%E8%BB%A2%E3%81%AE%E5%8E%9F%E5%89%87)

[依存関係逆転の原則について · SunriseDigital/work-shop Wiki](https://github.com/SunriseDigital/work-shop/wiki/%E4%BE%9D%E5%AD%98%E9%96%A2%E4%BF%82%E9%80%86%E8%BB%A2%E3%81%AE%E5%8E%9F%E5%89%87%E3%81%AB%E3%81%A4%E3%81%84%E3%81%A6)


### 具体例
コードを使用して具体例を示す。コードはGoで記述する。Goにこれまで馴染みのない方もなんとなくコードを見ればわかるかと思う。

これは後ほど記述するクリーンアーキテクチャで記述したコードの一部を切り取ったものだ。
クリーンアーキテクチャやコード全体は後述する。
この例では、ユースケースであるProgrammingLangUseCaseから使用されるデータベース周りの具体的な操作を行う構造体に焦点を当てる。

ProgrammingLangUseCaseから使用され、実際に操作を行うのはProgrammingLangDAOだが、ProgrammingLangUseCaseは、ProgrammingLangDAOをそのままProgrammingLangDAOとしては使用していない。(UseCaseやRepositoryについて、詳しくは[クリーンアーキテクチャ(The Clean Architecture翻訳) | blog.tai2.net](https://blog.tai2.net/the_clean_architecture.html)を参照)
どうしているかというと、ProgrammingLangRepositoryというインターフェースを定義し、その実装としてProgrammingLangDAOを使用している。
ProgrammingLangUseCaseは、ProgrammingLangRepositoryは知っているが、ProgrammingLangDAOは知らない。

なので、その部分はProgrammingLangRepositoryを実装している構造体ならば、何にでも差し替えることができる。
例えば、今回は、ProgrammingLangDAOはRDB(MySQL)の操作を実装しているが、ProgrammingLangRepositoryのインターフェースを満たしたNoSQLを操作する構造体に差し替えることもできるかもしれないし、メモリに保存する構造体に差し替えすることもできる。また、単体テストの際に、モックに差し替えることができる。これは単体テストを行う際には大きなメリットとなる。(単体テストについては後述する)

クラス図ぽいものを描くと以下のようなものになる。


![CleanArch.png](https://qiita-image-store.s3.amazonaws.com/0/145611/46c0f51b-5fa4-ad15-e33d-cbf9b38c3945.png)


上記のUMLのようにProgrammingLangUseCase(上位のレイヤー)がProgrammingLangDAOや、MockProgrammingLangRepository(下位レイヤー)に直接依存するのではなく、ProgrammingLangRepository(下位レイヤーの抽象)に依存し、ProgrammingLangDAOや、MockProgrammingLangRepository(下位レイヤー)は、ProgrammingLangRepository(下位レイヤーの抽象)の実装のため、下位レイヤーから下位レイヤーの抽象へ矢印が逆向きになるため、依存関係逆転の法則というらしい。

!注意1 : なんとなくUML図ぽく買いたものである。(厳密なUML図ではない)<br>
!注意2 : 実際のコードにはもう少しメソッドがあるが、説明のためだけの図なので、図には書かない。<br>

#### ProgrammingLangUseCase(上位レイヤ)
ProgrammingLangUseCase は、ProgrammingLangRepositoryを通して、ProgrammingLangDAOやMockProgrammingLangRepositoryを使用する。そのため、ProgrammingLangUseCase は、直接的には、具象であるProgrammingLangDAOやMockProgrammingLangRepositoryを知らない。

```go:program_lang_usecase.go

package usecase

import (
	"context"
	"time"

	"github.com/SekiguchiKai/clean-architecture-with-go/server/domain/model"
	"github.com/SekiguchiKai/clean-architecture-with-go/server/domain/repository"
	"github.com/SekiguchiKai/clean-architecture-with-go/server/usecase/input"
	"github.com/pkg/errors"
)

// ProgrammingLangUseCase は、ProgrammingLangのUseCase。
type ProgrammingLangUseCase struct {
	Repo repository.ProgrammingLangRepository
}

// NewProgrammingLangUseCase は、ProgrammingLangUseCaseを生成し、返す。
func NewProgrammingLangUseCase(repo repository.ProgrammingLangRepository) input.ProgrammingLangInputPort {
	return &ProgrammingLangUseCase{
		Repo: repo,
	}
}

// List は、ProgrammingLangの一覧を返す。
func (u *ProgrammingLangUseCase) List(ctx context.Context, limit int) ([]*model.ProgrammingLang, error) {
	limit = ManageLimit(limit, MaxLimit, MinLimit, DefaultLimit)
	return u.Repo.List(ctx, limit)
}

// Get は、ProgrammingLang1件返す。
func (u *ProgrammingLangUseCase) Get(ctx context.Context, id int) (*model.ProgrammingLang, error) {
	return u.Repo.Read(ctx, id)
}

// Create は、ProgrammingLangを生成する。
func (u *ProgrammingLangUseCase) Create(ctx context.Context, param *model.ProgrammingLang) (*model.ProgrammingLang, error) {
	lang, err := u.Repo.ReadByName(ctx, param.Name)
	if lang != nil {
		return nil, &model.AlreadyExistError{
			ID:        lang.ID,
			Name:      lang.Name,
			ModelName: model.ModelNameProgrammingLang,
		}
	}

	if _, ok := errors.Cause(err).(*model.NoSuchDataError); !ok {
		return nil, errors.WithStack(err)
	}

	param.CreatedAt = time.Now().UTC()
	param.UpdatedAt = time.Now().UTC()

	lang, err = u.Repo.Create(ctx, param)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return lang, nil
}

// Update は、ProgrammingLangを更新する。
func (u *ProgrammingLangUseCase) Update(ctx context.Context, id int, param *model.ProgrammingLang) (*model.ProgrammingLang, error) {
	lang, err := u.Repo.Read(ctx, id)
	if lang == nil {
		return nil, &model.NoSuchDataError{
			ID:        id,
			Name:      param.Name,
			ModelName: model.ModelNameProgrammingLang,
		}
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	lang.ID = id
	lang.Name = param.Name
	lang.Feature = param.Feature
	lang.UpdatedAt = time.Now().UTC()

	return u.Repo.Update(ctx, lang)
}

// Delete は、ProgrammingLangを削除する。
func (u *ProgrammingLangUseCase) Delete(ctx context.Context, id int) error {
	return u.Repo.Delete(ctx, id)
}

```

#### ProgrammingLangRepository(インターフェース部分)
ここでは、実際のデータベースの操作のインターフェースを定義している。個々のデータベースの操作(例えば、MySQLやPostgreSQL、あるいはそれを模したモックなど)という具体的なことに対して、ここで定義しているのは、データベースの操作をまとめた抽象的なものであることに注目して欲しい。これは、具体的なものが共通でもつ変わりにくい抽象的な部分をインターフェースで表したものだ。

```go:program_lang_repository.go
package repository

import (
	"context"

	"github.com/SekiguchiKai/clean-architecture-with-go/server/domain/model"
)

// ProgrammingLangRepository は、ProgrammingLangのRepository。
type ProgrammingLangRepository interface {
	List(ctx context.Context, limit int) ([]*model.ProgrammingLang, error)
	Create(ctx context.Context, lang *model.ProgrammingLang) (*model.ProgrammingLang, error)
	Read(ctx context.Context, id int) (*model.ProgrammingLang, error)
	ReadByName(ctx context.Context, name string) (*model.ProgrammingLang, error)
	Update(ctx context.Context, lang *model.ProgrammingLang) (*model.ProgrammingLang, error)
	Delete(ctx context.Context, id int) error
}
```

#### ProgrammingLangDAO(データベース操作実装部分)
具体的なSQL型のデータベースの操作を行う構造体(Javaとかでいうところのクラスみたいなもの)。
`ProgrammingLangRepository` で定義した各メソッドを実装している。ProgrammingLangDAOは、ProgrammingLangRepositoryを満たしているので、ProgrammingLangRepositoryとして使用することができる。


```go:program_lang_dao.go
package rdb

import (
	"context"
	"fmt"

	"github.com/SekiguchiKai/clean-architecture-with-go/server/domain/model"
	"github.com/SekiguchiKai/clean-architecture-with-go/server/domain/repository"
	"github.com/pkg/errors"
)

// ProgrammingLangDAO は、ProgrammingLangのDAO。
type ProgrammingLangDAO struct {
	SQLManager SQLManagerInterface
}

// NewProgrammingLangDAO は、ProgrammingLangDAO生成して返す。
func NewProgrammingLangDAO(manager SQLManagerInterface) repository.ProgrammingLangRepository {
	fmt.Printf("NewProgrammingLangDAO")

	return &ProgrammingLangDAO{
		SQLManager: manager,
	}
}

// ErrorMsg は、エラー文を生成し、返す。
func (dao *ProgrammingLangDAO) ErrorMsg(method string, err error) error {
	return &model.DBError{
		ModelName: model.ModelNameProgrammingLang,
		DBMethod:  method,
		Detail:    err.Error(),
	}
}

// Create は、レコードを1件生成する。
func (dao *ProgrammingLangDAO) Create(ctx context.Context, lang *model.ProgrammingLang) (*model.ProgrammingLang, error) {
	query := "INSERT INTO programming_languages (name, feature, created_at, updated_at) VALUES (?, ?, ?, ?)"
	stmt, err := dao.SQLManager.PrepareContext(ctx, query)
	if err != nil {
		return nil, dao.ErrorMsg(model.DBMethodCreate, err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, lang.Name, lang.Feature, lang.CreatedAt, lang.UpdatedAt)
	if err != nil {
		return nil, dao.ErrorMsg(model.DBMethodCreate, err)
	}

	affect, err := result.RowsAffected()
	if affect != 1 {
		err = fmt.Errorf("%s: %d ", TotalAffected, affect)
		return nil, dao.ErrorMsg(model.DBMethodUpdate, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, dao.ErrorMsg(model.DBMethodCreate, err)
	}

	lang.ID = int(id)

	return lang, nil
}

// List は、レコードの一覧を取得して返す。
func (dao *ProgrammingLangDAO) List(ctx context.Context, limit int) ([]*model.ProgrammingLang, error) {
	query := "SELECT id, name, feature, created_at, updated_at FROM programming_languages ORDER BY name LIMIT ?"
	return dao.list(ctx, query, limit)
}

// list は、レコードの一覧を取得して返す。
func (dao *ProgrammingLangDAO) list(ctx context.Context, query string, args ...interface{}) ([]*model.ProgrammingLang, error) {
	stmt, err := dao.SQLManager.PrepareContext(ctx, query)
	if err != nil {
		return nil, dao.ErrorMsg(model.DBMethodList, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, dao.ErrorMsg(model.DBMethodList, err)
	}
	defer rows.Close()

	langSlice := make([]*model.ProgrammingLang, 0)
	for rows.Next() {
		lang := &model.ProgrammingLang{}

		err = rows.Scan(
			&lang.ID,
			&lang.Name,
			&lang.Feature,
			&lang.CreatedAt,
			&lang.UpdatedAt,
		)

		if err != nil {
			return nil, dao.ErrorMsg(model.DBMethodList, err)
		}
		langSlice = append(langSlice, lang)
	}

	return langSlice, nil
}

// Read は、レコードを1件取得して返す。
func (dao *ProgrammingLangDAO) Read(ctx context.Context, id int) (*model.ProgrammingLang, error) {
	query := "SELECT id, name, feature, created_at, updated_at FROM programming_languages WHERE ID=?"

	stmt, err := dao.SQLManager.PrepareContext(ctx, query)
	if err != nil {
		return nil, dao.ErrorMsg(model.DBMethodRead, err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	lang := &model.ProgrammingLang{}

	err = row.Scan(
		&lang.ID,
		&lang.Name,
		&lang.Feature,
		&lang.CreatedAt,
		&lang.UpdatedAt,
	)

	if err != nil {
		return nil, dao.ErrorMsg(model.DBMethodRead, err)
	}

	return lang, nil
}

// ReadByName は、指定したNameを保持するレコードを1返す。
func (dao *ProgrammingLangDAO) ReadByName(ctx context.Context, name string) (*model.ProgrammingLang, error) {
	query := "SELECT id, name, feature, created_at, updated_at FROM programming_languages WHERE name=? ORDER BY name LIMIT ?"
	langSlice, err := dao.list(ctx, query, name, 1)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(langSlice) == 0 {
		return nil, &model.NoSuchDataError{
			Name:      name,
			ModelName: model.ModelNameProgrammingLang,
		}
	}

	return langSlice[0], nil
}

// Update は、レコードを1件更新する。
func (dao *ProgrammingLangDAO) Update(ctx context.Context, lang *model.ProgrammingLang) (*model.ProgrammingLang, error) {
	query := "UPDATE programming_languages SET name=?, feature=?, created_at=?, updated_at=? WHERE id=?"

	stmt, err := dao.SQLManager.PrepareContext(ctx, query)
	defer stmt.Close()

	if err != nil {
		return nil, dao.ErrorMsg(model.DBMethodUpdate, err)
	}

	result, err := stmt.ExecContext(ctx, lang.Name, lang.Feature, lang.CreatedAt, lang.UpdatedAt, lang.ID)
	if err != nil {
		return nil, dao.ErrorMsg(model.DBMethodUpdate, err)
	}

	affect, err := result.RowsAffected()
	if affect != 1 {
		err = fmt.Errorf("%s: %d ", TotalAffected, affect)
		return nil, dao.ErrorMsg(model.DBMethodUpdate, err)
	}

	return lang, nil
}

// Delete は、レコードを1件削除する。
func (dao *ProgrammingLangDAO) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM programming_languages WHERE id=?"

	stmt, err := dao.SQLManager.PrepareContext(ctx, query)
	if err != nil {
		return dao.ErrorMsg(model.DBMethodDelete, err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return dao.ErrorMsg(model.DBMethodDelete, err)
	}

	affect, err := result.RowsAffected()
	if err != nil {
		return dao.ErrorMsg(model.DBMethodDelete, err)
	}
	if affect != 1 {
		err = fmt.Errorf("%s: %d ", TotalAffected, affect)
		return dao.ErrorMsg(model.DBMethodDelete, err)
	}

	return nil
}

```



#### MockProgrammingLangRepository(モック)
データベースの操作を模したモック。
[gomock](https://github.com/golang/mock)というモック生成用のツールで自動生成している。
モックの構造体もProgrammingLangRepositoryを満たしているので、ProgrammingLangRepositoryとして使用することができる。実際にProgrammingLangRepository(の実装)を使用するレイヤーのテストをする際には、ProgrammingLangRepositoryの実装としてProgrammingLangDAOではなく、このモックを使用する。

```go:program_lang_repository_mock.go

// Code generated by MockGen. DO NOT EDIT.
// Source: domain/repository/programming_lang_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	model "github.com/SekiguchiKai/clean-architecture-with-go/server/domain/model"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockProgrammingLangRepository is a mock of ProgrammingLangRepository interface
type MockProgrammingLangRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProgrammingLangRepositoryMockRecorder
}

// MockProgrammingLangRepositoryMockRecorder is the mock recorder for MockProgrammingLangRepository
type MockProgrammingLangRepositoryMockRecorder struct {
	mock *MockProgrammingLangRepository
}

// NewMockProgrammingLangRepository creates a new mock instance
func NewMockProgrammingLangRepository(ctrl *gomock.Controller) *MockProgrammingLangRepository {
	mock := &MockProgrammingLangRepository{ctrl: ctrl}
	mock.recorder = &MockProgrammingLangRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProgrammingLangRepository) EXPECT() *MockProgrammingLangRepositoryMockRecorder {
	return m.recorder
}

// List mocks base method
func (m *MockProgrammingLangRepository) List(ctx context.Context, limit int) ([]*model.ProgrammingLang, error) {
	ret := m.ctrl.Call(m, "List", ctx, limit)
	ret0, _ := ret[0].([]*model.ProgrammingLang)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockProgrammingLangRepositoryMockRecorder) List(ctx, limit interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockProgrammingLangRepository)(nil).List), ctx, limit)
}

// Create mocks base method
func (m *MockProgrammingLangRepository) Create(ctx context.Context, lang *model.ProgrammingLang) (*model.ProgrammingLang, error) {
	ret := m.ctrl.Call(m, "Create", ctx, lang)
	ret0, _ := ret[0].(*model.ProgrammingLang)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockProgrammingLangRepositoryMockRecorder) Create(ctx, lang interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockProgrammingLangRepository)(nil).Create), ctx, lang)
}

// Read mocks base method
func (m *MockProgrammingLangRepository) Read(ctx context.Context, id int) (*model.ProgrammingLang, error) {
	ret := m.ctrl.Call(m, "Read", ctx, id)
	ret0, _ := ret[0].(*model.ProgrammingLang)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read
func (mr *MockProgrammingLangRepositoryMockRecorder) Read(ctx, id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockProgrammingLangRepository)(nil).Read), ctx, id)
}

// ReadByName mocks base method
func (m *MockProgrammingLangRepository) ReadByName(ctx context.Context, name string) (*model.ProgrammingLang, error) {
	ret := m.ctrl.Call(m, "ReadByName", ctx, name)
	ret0, _ := ret[0].(*model.ProgrammingLang)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadByName indicates an expected call of ReadByName
func (mr *MockProgrammingLangRepositoryMockRecorder) ReadByName(ctx, name interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadByName", reflect.TypeOf((*MockProgrammingLangRepository)(nil).ReadByName), ctx, name)
}

// Update mocks base method
func (m *MockProgrammingLangRepository) Update(ctx context.Context, lang *model.ProgrammingLang) (*model.ProgrammingLang, error) {
	ret := m.ctrl.Call(m, "Update", ctx, lang)
	ret0, _ := ret[0].(*model.ProgrammingLang)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockProgrammingLangRepositoryMockRecorder) Update(ctx, lang interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockProgrammingLangRepository)(nil).Update), ctx, lang)
}

// Delete mocks base method
func (m *MockProgrammingLangRepository) Delete(ctx context.Context, id int) error {
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockProgrammingLangRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockProgrammingLangRepository)(nil).Delete), ctx, id)
}

```

# 単体テスト編
このセクションは、以下の2つの記事を大変参考にさせていただいた。
[mercari.go #1 で「もう一度テストパターンを整理しよう」というタイトルで登壇しました - アルパカ三銃士](https://codehex.hateblo.jp/entry/2018/07/03/211839)

[Goにおけるテスト可能な設計](https://www.slideshare.net/shogoosawa581/go-77254684)

### そもそも単体テストとは何かということを振り返る
単体テストについての説明は色々とあると思うが機能テストと比較して書かれた以下の説明がわかりやすい。

> Unit test(単体テスト)
　・単一の関数やメソッドなどをテスト
> Functional test(機能テスト)
　・リクエストからレスポンスまでのテスト

引用元 : [mercari.go #1 で「もう一度テストパターンを整理しよう」というタイトルで登壇しました - アルパカ三銃士](https://codehex.hateblo.jp/entry/2018/07/03/211839)

### テストブル
よく、テストダブルという言葉を聞いたことはないだろうか。
テストダブルとは、
> ソフトウェアテストにおいて、テスト対象が依存しているコンポーネントを置き換える代用品のこと。ダブルは代役、影武者を意味する。
テストを実行するには、テスト対象のシステム (SUT; System Under Test) に加えて、テスト対象が依存するコンポーネント (DOC; Depend-On Component) が必要になる。しかし、依存コンポーネントは、常に利用できるとは限らない。

> こういった問題を回避するには、依存コンポーネントを、テスト用のコンポーネントと入れ替えるテクニックが利用できる。この代用のコンポーネントを、テストダブルと呼ぶ。

引用元 : [テストダブル - Wikipedia](https://ja.wikipedia.org/wiki/%E3%83%86%E3%82%B9%E3%83%88%E3%83%80%E3%83%96%E3%83%AB)

要するに、あるコンポーネントをテストする際に、そのテスト対象のコンポーネントが依存しているコンポーネントが利用できなかったりするので、それをテスト用に作ったものに置きかえようねという話。

これの総称がテストダブルで、その具体的な方法にはモックやスタブなどがある。
各々の違いについては、[TDD > モック / スタブ - Qiita](https://qiita.com/7of9/items/8e2cb2070f2b2ea4e5ec)等で確認いただきたい。

今回は、その中でモックを使用する。

## 単体テストでインターフェースをうまく利用する
先ほど、引用で単体テストは「単一の関数やメソッドなどをテスト」するということがわかった。
A=>B=>Cという依存関係がコードがあるとする。AはBに依存し、BはCに依存するとする。この場合、Aのテストを行おうとすると、BやCまで呼び出す必要が出て来てしまう。
先ほどの単体テストの定義だと、Aの単体テストはAのみをテストするものなはずなのに、A以外のBやCも利用することになってしまう。
これは真の意味で単体テストと言えるのだろうか...

依存関係のある中で、単体テストをうまく行うのにインターフェイスとポリモーフィズムを使用するとAの単体テストを行うのに、実際のBやCを利用しなくてもよくなる。

実際の例は、先ほどのDIP(依存関係逆転の法則)のセクションで示したものを参照いただきたい。原理としては、Aの単体テストをする際に、依存しているBやCをそのまま使うのではなく、Bをモックに入れ替えている。
これは、AからBを利用する際に、Bの具象クラスをそのまま利用するのではなくて、Bの具象クラスがその実装となるインターフェイスを定義して、それをAは利用しているからなせる技だ。

具体的にいうとProgrammingLangRepositoryというインターフェースを定義し、製品コードではこのProgrammingLangRepositoryの実装であるProgrammingLangDAOを使用してDBの操作を行い、ProgrammingLangUseCaseのテストでは、ProgrammingLangRepository実装であるMockProgrammingLangRepositoryに差し替えているのだ。
モックもインターフェースを実装した具象クラスの1つであるというわけだ。

**ProgrammingLangDAO is a ProgrammingLanRepository** であり、<br>
**MockProgrammingLangRepository is a ProgrammingLanRepository** でもあるという事だ。

# クリーンアーキテクチャ編
変更に強く、テストがしやすいということで(もちろん他にも利点はたくさんある)最近何かと話題に上がることの多いクリーンアーキテクチャ。
これの何が優れているのかということをこれまでの説明から紐解きたい。
ただし、既にクリーンアーキテクチャの優れた部分は他の記事等でも紹介されているので、ここでは、これまでの記事の内容に沿ったものだけに焦点を当てたい。

このセクションでは以下の記事を非常に参考にさせていただいた。
[Clean ArchitectureでAPI Serverを構築してみる - Qiita](https://qiita.com/hirotakan/items/698c1f5773a3cca6193e)

[Goでクリーンアーキテクチャを試す | POSTD](https://postd.cc/golang-clean-archithecture/)

[Goのサーバサイド実装におけるレイヤ設計とレイヤ内実装について考える](https://www.slideshare.net/pospome/go-80591000)

[クリーンアーキテクチャ(The Clean Architecture翻訳) | blog.tai2.net](https://blog.tai2.net/the_clean_architecture.html)

## 変更に強くなる編に合致する点
### 依存の方向性
> 内側に向かってのみ依存することができる。
というように循環依存しないようにする。

詳しくは、[クリーンアーキテクチャ(The Clean Architecture翻訳) | blog.tai2.net](https://blog.tai2.net/the_clean_architecture.html)を参照。

### レイヤーとDIP
詳しくは[クリーンアーキテクチャ(The Clean Architecture翻訳) | blog.tai2.net](https://blog.tai2.net/the_clean_architecture.html)を参照いただきたいが、レイヤーを分けて、レイヤ間の境界をまたがるときには、疎結合になるようにDIPを用いることが多い。これを行うことで、あるレイヤのコードが変更になったときに、別のレイヤーに影響を及ぼしにくい。

## 単体テスト編に合致する点
### モックにできる
レイヤを分けて、レイヤ間の境界をまたがるときには、疎結合になるようにDIPを用いるので、依存している他のレイヤはモックに差し替えることができるため、単体テストがしやすい。

## 実際のコード
実際にコードを書いてみた。
https://github.com/SekiguchiKai/clean-architecture-with-go


## 参考文献
エリック・エヴァンス(著)、 今関 剛 (監修)、 和智 右桂  (翻訳) (2011/4/9)『エリック・エヴァンスのドメイン駆動設計 (IT Architects’Archive ソフトウェア開発の実践)』 翔泳社

Robert C.Martin (著)、 角 征典  (翻訳)、 高木 正弘 (翻訳)　(2018/7/27)『Clean Architecture 達人に学ぶソフトウェアの構造と設計』 KADOKAWA

アラン・シャロウェイ (著)、 ジェームズ・R・トロット (著)、 村上 雅章 (翻訳) (2014/3/11)『オブジェクト指向のこころ (SOFTWARE PATTERNS SERIES)』 丸善出版

結城 浩 (2004/6/19)『増補改訂版Java言語で学ぶデザインパターン入門』 ソフトバンククリエイティブ

 InfoQ.com、徳武 聡(翻訳) (2009年6月7日) 『Domain Driven Design（ドメイン駆動設計） Quickly 日本語版』 InfoQ.com
[Domain Driven Design（ドメイン駆動設計） Quickly 日本語版](https://www.infoq.com/jp/minibooks/domain-driven-design-quickly)

中山 清喬、国本 大悟 (2014/8/7)『スッキリわかるJava入門 第2版 スッキリわかるシリーズ』 インプレス

## 参考にさせていただいたサイト
### 変更に強くなる編
[実践DDD本の第4章「アーキテクチャ」 ～レイヤからヘキサゴナルへ～ (2/4)：CodeZine（コードジン）](https://codezine.jp/article/detail/9922?p=2)

[オブジェクト指向設計原則とは - Qiita](https://qiita.com/UWControl/items/98671f53120ae47ff93a)

[依存性逆転の原則 - Wikipedia](https://ja.wikipedia.org/wiki/%E4%BE%9D%E5%AD%98%E6%80%A7%E9%80%86%E8%BB%A2%E3%81%AE%E5%8E%9F%E5%89%87)

[依存関係逆転の原則について · SunriseDigital/work-shop Wiki](https://github.com/SunriseDigital/work-shop/wiki/%E4%BE%9D%E5%AD%98%E9%96%A2%E4%BF%82%E9%80%86%E8%BB%A2%E3%81%AE%E5%8E%9F%E5%89%87%E3%81%AB%E3%81%A4%E3%81%84%E3%81%A6)

[オブジェクト指向と10年戦ってわかったこと - Qiita](https://qiita.com/tutinoco/items/6952b01e5fc38914ec4e#%E3%83%9D%E3%83%AA%E3%83%A2%E3%83%BC%E3%83%95%E3%82%A3%E3%82%BA%E3%83%A0)

[15分でわかる かんたんオブジェクト指向 - Qiita](https://qiita.com/koher/items/6878c80014992900add7#%E3%83%9D%E3%83%AA%E3%83%A2%E3%83%BC%E3%83%95%E3%82%A3%E3%82%BA%E3%83%A0)

参考 : [モジュール結合度とは - IT用語辞典 Weblio辞書](https://www.weblio.jp/content/%E3%83%A2%E3%82%B8%E3%83%A5%E3%83%BC%E3%83%AB%E7%B5%90%E5%90%88%E5%BA%A6)

### 単体テスト編
[mercari.go #1 で「もう一度テストパターンを整理しよう」というタイトルで登壇しました - アルパカ三銃士](https://codehex.hateblo.jp/entry/2018/07/03/211839)

[Goにおけるテスト可能な設計](https://www.slideshare.net/shogoosawa581/go-77254684)

[テストダブル - Wikipedia](https://ja.wikipedia.org/wiki/%E3%83%86%E3%82%B9%E3%83%88%E3%83%80%E3%83%96%E3%83%AB)

[TDD > モック / スタブ - Qiita](https://qiita.com/7of9/items/8e2cb2070f2b2ea4e5ec)等で確認いただきたい。

### クリーンアーキテクチャ 編
[The Clean Architecture | 8th Light](https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html)

[クリーンアーキテクチャ(The Clean Architecture翻訳) | blog.tai2.net](https://blog.tai2.net/the_clean_architecture.html)

[Clean Architecture │ nrslib](https://nrslib.com/clean-architecture/)

[Goのサーバサイド実装におけるレイヤ設計とレイヤ内実装について考える](https://www.slideshare.net/pospome/go-80591000)

[Clean ArchitectureでAPI Serverを構築してみる - Qiita](https://qiita.com/hirotakan/items/698c1f5773a3cca6193e)

[Goでクリーンアーキテクチャを試す | POSTD](https://postd.cc/golang-clean-archithecture/)

[持続可能な開発を目指す ~ ドメイン・ユースケース駆動（クリーンアーキテクチャ） + 単方向に制限した処理 + FRP - Qiita](https://qiita.com/kondei/items/41c28674c1bfd4156186)
