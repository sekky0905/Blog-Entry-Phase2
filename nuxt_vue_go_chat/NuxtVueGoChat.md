Nuxt.js(Vue.js)とGoでSPA + API(レイヤードアーキテクチャ)でチャットアプリを実装してみた

# 概要

Nuxt.js(Vue.js)とレイヤードアーキテクチャのお勉強のために簡単なチャットアプリを実装してみた。
SPA + APIと言った形になっている。

# 機能

機能はだいたい以下のような感じ。

- ログイン機能
- サインアップ機能
- スレッド一覧表示機能
- スレッド作成機能
    - ログインしたユーザーは誰でもスレッドを作成できること
- コメント一覧表示機能
    - スレッドをクリックすると、そのスレッド内のコメント一覧が表示されること
- スレッド内でのコメント作成機能
    - ログインしたユーザーは誰でもどのスレッド内でもコメントできること
- スレッド内でのコメント削除機能
    - 自分のコメントのみ削除できること
- ログアウト機能

# コード
コード全体は[ここ](https://github.com/sekky0905/nuxt-vue-go-chat)。

# 技術

## サーバーサイド

### アーキテクチャ

DDD本に出てくるレイヤードアーキテクチャをベースに以下の書籍や記事を参考にさせていただき実装した。

- [Goを運用アプリケーションに導入する際のレイヤ構造模索の旅路 | Go Conference 2018 Autumn 発表レポート - BASE開発チームブログ](https://devblog.thebase.in/entry/2018/11/26/102401)
- [GoでのAPI開発現場のアーキテクチャ実装事例 / go-api-architecture-practical-example - Speaker Deck](https://speakerdeck.com/hgsgtk/go-api-architecture-practical-example)
- [ボトムアップドメイン駆動設計 │ nrslib](https://nrslib.com/bottomup-ddd/)
- エリック・エヴァンス(著)、今関 剛 (監修)、和智 右桂 (翻訳) (2011/4/9)『エリック・エヴァンスのドメイン駆動設計 (IT Architects’Archive ソフトウェア開発の実践)』 翔泳社
- pospome『pospomeのサーバサイドアーキテクチャ

実際のpackage構成は以下のような感じ。

```
├── interface
│   └── controller // サーバへの入力と出力を扱う責務。
├── application // 薄く保ち、やるべき作業の調整を行う責務。
├── domain
│   ├── model // ビジネスの概念とビジネスロジック。
│   ├── service // EntityでもValue Objectでもないドメイン層のロジック。
│   └── repository // infra/dbへのポート。
├── infra // 技術的なものの提供
│    ├── db // DBの技術に関すること。
│    ├── logger // Logの技術に関すること。
│    └── router // Routingの技術に関すること。 
├── middleware // リクエスト毎に差し込む処理をまとめたミドルウェア
├── util 
└── testutil
```

上記のpackage以外に `application/mock`、`domain/service/mock` 、`infra/db/mock` には、mockという格納する用のpackageもあり、そこに各々のレイヤーのmock用のファイルを置いている。(詳しくは後述)

#### 依存関係

依存関係としてはざっくり、`interface/controller` → `application` → `dmain/repository` or `dmain/service` ← `infra/db` という形になっている。

参考: [GoでのAPI開発現場のアーキテクチャ実装事例 / go-api-architecture-practical-example - Speaker Deck](https://speakerdeck.com/hgsgtk/go-api-architecture-practical-example?slide=16)

`domain/~` と `infra/db` で矢印が逆になっているのは、依存関係が逆転しているため。
詳しくは [その設計、変更に強いですか?単体テストできますか?...そしてクリーンアーキテクチャ - Qiita](https://qiita.com/Sekky0905/items/2436d669ff5d4491c527)を参照。

先ほどの矢印の中で、`domain/model` は記述しなかったが、 `domain/model` は、`interface/controller` や `application` 等からも依存されている。純粋なレイヤードアーキテクチャでは、各々のレイヤーは自分の下のレイヤーにのみ依存するといったものがあるかもしれないが、それを実現するためにDTO等を用意するのが大変だったから。

#### 各レイヤーでのinterfaceの定義とテスト

`applicaion` 、 `/domain/service` 、`infra/db` (定義先は、`/domain/repository` ) には、 `interface` を定義し、他のレイヤーからはその `interface` に依存させるようにしている。こうするとこれらを使用する側は、抽象に依存するようになるので、抽象を実装する具象を変化させても、使用する側(依存する側)は、その影響を受けにくい。

実際に各レイヤーを使用する側のレイヤのテストの際には、使用されるレイヤーを実際のコードではなく、Mock用のものに差し替えている。各々のレイヤーに存在する `mock` というpackageにmock用のコードを置いている。このモック用のコードは、[gomock](https://github.com/golang/mock)を使用して、自動生成している。

この辺のことについては、
[その設計、変更に強いですか?単体テストできますか?...そしてクリーンアーキテクチャ - Qiita](https://qiita.com/Sekky0905/items/2436d669ff5d4491c527#%E3%82%A4%E3%83%B3%E3%82%BF%E3%83%BC%E3%83%95%E3%82%A7%E3%83%BC%E3%82%B9%E3%81%A8%E3%83%9D%E3%83%AA%E3%83%A2%E3%83%BC%E3%83%95%E3%82%A3%E3%82%BA%E3%83%A0) という記事を以前書いたので、詳しくはこちらを参照いただきたい。

### エラーハンドリング

エラーハンドリングは以下のように行なっている。

- 以下のような形で `errors.Wrap` を使用してオリジナルのエラーを包む

```go
if err := Hoge(); err != nil {
    return errors.Wrap(オリジナルエラー, '状況の説明'
}
```

- 独自のエラー型を定義している
- エラーは基本的に各々のレイヤーで握りつぶさず、`interface/controller` レイヤーまで伝播させる
- 最終的には、`interface/controller` でエラーの型によって、レスポンスとして返すメッセージやステータスコードを選択する

参考
[Golangのエラー処理とpkg/errors | SOTA](https://deeeet.com/writing/2016/04/25/go-pkg-errors/)

### ログイン周り

- 外部サービスを使用せず、自前で簡単なものを実装した。
- パスワードは [bcrypt](https://godoc.org/golang.org/x/crypto/bcrypt)を使用した
    - 参考[【Go言語】パスワードをハッシュ化(bcrypt) - まったり技術ブログ](https://blog.motikan2010.com/entry/2017/02/13/%E3%80%90Go%E8%A8%80%E8%AA%9E%E3%80%91%E3%83%91%E3%82%B9%E3%83%AF%E3%83%BC%E3%83%89%E3%82%92%E3%83%8F%E3%83%83%E3%82%B7%E3%83%A5%E5%8C%96%28bcrypt%29)
- 普通にCookieとSessionを使用した
- ログインが必要なAPIには `gin` の `middleware` を使用して、ログイン済みでないクライアントからのリクエストは `401 Unauthorized` を返すようにした
    - [gin middleware](https://github.com/gin-gonic/gin#using-middleware)
    - [GolangのWebフレームワークginのmiddlewareについての覚書 - Qiita](https://qiita.com/tobita0000/items/d2309cc3f0dd32006ead)

### DB周り

- MySQLを使用した
- DBテスト部分は、DBサーバを立てたわけではなく、[DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)を使用し、モックで行なった
    - 以下を使用してDBサーバーを立てて行うのも良いかも
        - [ory/dockertest](https://github.com/ory/dockertest)
            - Dockerを使う場合
        - [lestrrat-go/test-mysqld](https://github.com/lestrrat-go/test-mysqld)
            - Dockerを使わない場合
- DB操作周りの実装に関しては、[database/sql](https://golang.org/pkg/database/sql/)packageをそのまま使用し、ORMやその他のライブラリは特に使用していない
- トランザクションは、`application` レイヤでかける
- 以下のようなSQL周りの `interface` を作成
    - 参考: [mercari.go #1 で「もう一度テストパターンを整理しよう」というタイトルで登壇しました - アルパカ三銃士](https://codehex.hateblo.jp/entry/2018/07/03/211839)

```go
package query

import (
	"context"
	"database/sql"
)

// DBManager is the manager of SQL.
type DBManager interface {
	SQLManager
	Beginner
}

// TxManager is the manager of Tx.
type TxManager interface {
	SQLManager
	Commit() error
	Rollback() error
}

// SQLManager is the manager of DB.
type SQLManager interface {
	Querier
	Preparer
	Executor
}

type (
	// Executor is interface of Execute.
	Executor interface {
		Exec(query string, args ...interface{}) (sql.Result, error)
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	}

	// Preparer is interface of Prepare.
	Preparer interface {
		Prepare(query string) (*sql.Stmt, error)
		PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	}

	// Querier is interface of Query.
	Querier interface {
		Query(query string, args ...interface{}) (*sql.Rows, error)
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	}

	// Beginner is interface of Begin.
	Beginner interface {
		Begin() (TxManager, error)
	}
)

```

- `application` レイヤーでは以下のようにフィールで `query.DBManager` を所持する
    - そうすることで `SQLManager` も `Begin()` し `TxManager` もどちらも `application` レイヤーで扱うことができる( `application` レイヤで直接使用するわけではなく、 `domain/repository` に渡す)

```go
// threadService is application service of thread.
type threadService struct {
	m        query.DBManager
	service  service.ThreadService
	repo     repository.ThreadRepository
	txCloser CloseTransaction
}
```

- `domain/repository` の引数では `query.SQLManager` を受け取る
    - `query.TxManager` は、`query.SQLManager` も満たしているので、`query.TxManager` は、`query.SQLManager` のどちらも受け取ることができる

```go
// ThreadRepository is Repository of Thread.
type ThreadRepository interface {
	ListThreads(ctx context.Context, m query.SQLManager, cursor uint32, limit int) (*model.ThreadList, error)
	GetThreadByID(ctx context.Context, m query.SQLManager, id uint32) (*model.Thread, error)
	GetThreadByTitle(ctx context.Context, m query.SQLManager, name string) (*model.Thread, error)
	InsertThread(ctx context.Context, m query.SQLManager, thead *model.Thread) (uint32, error)
	UpdateThread(ctx context.Context, m query.SQLManager, id uint32, thead *model.Thread) error
	DeleteThread(ctx context.Context, m query.SQLManager, id uint32) error
}
```

- 以下のようなRollbackやCommitを行う関数を作成しておく

```go
// CloseTransaction executes post process of tx.
func CloseTransaction(tx query.TxManager, err error) error {
	if p := recover(); p != nil { // rewrite panic
		err = tx.Rollback()
		err = errors.Wrap(err, "failed to roll back")
		panic(p)
	} else if err != nil {
		err = tx.Rollback()
		err = errors.Wrap(err, "failed to roll back")
	} else {
		err = tx.Commit()
		err = errors.Wrap(err, "failed to commit")
	}
	return err
}
```

- `application` レイヤでは、`defer` で `CloseTransaction` を呼び出す(ここでは `a.txCloser` になっている)

```go
// CreateThread creates Thread.
func (a *threadService) CreateThread(ctx context.Context, param *model.Thread) (thread *model.Thread, err error) {
	tx, err := a.m.Begin()
	if err != nil {
		return nil, beginTxErrorMsg(err)
	}

	defer func() {
		if err := a.txCloser(tx, err); err != nil {
			err = errors.Wrap(err, "failed to close tx")
		}
	}()

	yes, err := a.service.IsAlreadyExistTitle(ctx, tx, param.Title)
	if yes {
		err = &model.AlreadyExistError{
			PropertyName:    model.TitleProperty,
			PropertyValue:   param.Title,
			DomainModelName: model.DomainModelNameThread,
		}
		return nil, errors.Wrap(err, "already exist id")
	}

	if _, ok := errors.Cause(err).(*model.NoSuchDataError); !ok {
		return nil, errors.Wrap(err, "failed is already exist id")
	}

	id, err := a.repo.InsertThread(ctx, tx, param)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert thread")
	}
	param.ID = id
	return param, nil
}
```

- 上記の処理ができるように `CloseTransaction` を `application` レイヤの構造体にDIしておく
    - Goでは関数もDIできる
        - [Goの構造体のフィールドに関数を埋め込んでDIを実現してみる - Qiita](https://qiita.com/Sekky0905/items/be5d83674a1aa78a397a)を参照

```go
// threadService is application service of thread.
type threadService struct {
	m        query.DBManager
	service  service.ThreadService
	repo     repository.ThreadRepository
	txCloser CloseTransaction
}
```

### 所感

- 依存関係が決まるのが良い
- 各レイヤが疎結合なので変更しやすく、テストもしやすいのは良い
- 各レイヤの責務がはっきり別れているので、どこに何を書けばいいかはわかりやすい
- コード量は増えるので、実装に時間がかかる
    - 決まったところは自動化できると良いかも
    - CRUDだけの小さなアプリケーションでは、大げさすぎるかもしれない

## フロントエンド

### アーキテクチャ

- 基本的には、Nuxt.jsのアーキテクチャに沿って実装を行なった。
- 状態管理に感じては、Vuexを使用した
    - 各々の `Component` 側( `pages` や `components` )からデータを使用したい場合には、Vuexを通じて使用した
    - データ、ロジックとビュー部分が綺麗に別れる

### 見た目

- [Vue.js](https://vuetifyjs.com/ja/)に全面的に乗っかった
- コメントの一覧部分のCSSは [CSSで作る！吹き出しデザインのサンプル19選](https://saruwakakun.com/html-css/reference/speech-bubble) を参考にさせていただいた

### 大きな流れ

大きな流れとしては、以下のような流れ。
`pasges` や `components` 等のビューでのイベントの発生 → `actions` でAPIへリクエスト → `mutations` で `state` 変更 → `pasges` や `components` 等のビューに反映される

他の流れもたくさんあるが、代表的なList処理とInput処理の流れを以下に記す。

#### List処理

- `pages` や `components` の `asyncData` 内で、`store.dispatch` を通じて、データ一覧を取得するアクション( `actions` )を呼び出す
- `store` の `actions` 内での処理を行う
    -  [axios](https://github.com/axios/axios)を使用してAPIにリクエストを送信する
    - APIから返却されたデータを引数に `mutations` を `commit` する。
- `mutations` での処理を行う
    - `state` を変更する
- `pages` や `components` のビューで取得したデータが表示される

#### Input処理

- `pages` や `components` で `stores` に定義した `action` や `state` を読み込んでおく
- `pages` や `components` の `data` 部分とformのinput部分等に `v-model` を使用して双方向データバインディングをしておく
    - [フォーム入力バインディング — Vue.js](https://jp.vuejs.org/v2/guide/forms.html)
- `pages` や `components` で表示しているビュー部分でイベントが生じる
    - form入力→submitなど
- sumitする時にクリックされるボタンに `@click=hoge` という形でイベントがそのElementで該当のイベントが生じた時に呼び出されるメソッド等を登録しておく
    - 上記の例では、 `click` イベントが生じると `hoge` メソッドが呼び出される
    - [イベントハンドリング — Vue.js](https://jp.vuejs.org/v2/guide/events.html)
- 呼び出されたメソッドの処理を行う
    - formのデータを元にデータを登録するアクション( `actions` )を呼び出す
- `store` の `actions` 内での処理を行う
    -  [axios](https://github.com/axios/axios)を使用してAPIにリクエストを送信する
    - APIから返却されたデータを引数に `mutations` を `commit` する。
- `mutations` での処理を行う
    - `state` を変更する
    - 登録した分のデータを一覧の `state` に追加する
- `pages` や `components` のビューで登録したデータが追加された一覧表示される

### 非同期部分

- `async/await` で処理
    - [async await の使い方 - Qiita](https://qiita.com/niusounds/items/37c1f9b021b62194e077)

### 所感

- Nuxt.jsを使用すると、レールに乗っかれて非常に楽
    - どこに何を実装すればいいか明白になるので迷わないで済む
    - 特にVuexを使用すると
        - データの流れが片方向になるのはわかりやすくて良い
        - ビュー、ロジック、データの責務がはっきりするのが良い
- Vuetifyを使用するとあまり凝らない画面であれば、短期間で実装できそう

# 参考文献

## サーバーサイド

- InfoQ.com、徳武 聡(翻訳) (2009年6月7日) 『Domain Driven Design（ドメイン駆動設計） Quickly 日本語版』 InfoQ.com
- エリック・エヴァンス(著)、今関 剛 (監修)、和智 右桂 (翻訳) (2011/4/9)『エリック・エヴァンスのドメイン駆動設計 (IT Architects’Archive ソフトウェア開発の実践)』 翔泳社
- pospome『pospomeのサーバサイドアーキテクチャ』

## フロントエンド
- 花谷拓磨 (2018/10/17)『Nuxt.jsビギナーズガイド』シーアンドアール研究所
- 川口 和也、喜多 啓介、野田 陽平、 手島 拓也、 片山 真也(2018/9/22)『Vue.js入門 基礎から実践アプリケーション開発まで』技術評論社


# 参考にさせていただいた記事

## サーバーサイド

- [Goを運用アプリケーションに導入する際のレイヤ構造模索の旅路 | Go Conference 2018 Autumn 発表レポート - BASE開発チームブログ](https://devblog.thebase.in/entry/2018/11/26/102401)

- [ボトムアップドメイン駆動設計 │ nrslib](https://nrslib.com/bottomup-ddd/)
- [GoでのAPI開発現場のアーキテクチャ実装事例 / go-api-architecture-practical-example - Speaker Deck](https://speakerdeck.com/hgsgtk/go-api-architecture-practical-example)
- [GoのAPIのテストにおける共通処理 – timakin – Medium](https://medium.com/@timakin/go-api-testing-173b97fb23ec)
- [mercari.go #1 で「もう一度テストパターンを整理しよう」というタイトルで登壇しました - アルパカ三銃士](https://codehex.hateblo.jp/entry/2018/07/03/211839)

## フロントエンド

- [Nuxt.js - ユニバーサル Vue.js アプリケーション](https://ja.nuxtjs.org/)
- [Vuex](https://vuex.vuejs.org/ja/)
- [Vuetify.js](https://vuetifyjs.com/ja/)
- [Vue.js](https://jp.vuejs.org/index.html)
- [CSSで作る！吹き出しデザインのサンプル19選](https://saruwakakun.com/html-css/reference/speech-bubble)

# 関連記事

- [その設計、変更に強いですか?単体テストできますか?...そしてクリーンアーキテクチャ - Qiita](https://qiita.com/Sekky0905/items/2436d669ff5d4491c527)
