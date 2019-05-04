MySQLのInnoDBのトランザクション周りについて調べたことをまとめた

元々NoSQLしか使用していなく、今年に入って転職して初めて業務でMySQLを使用しだした。
雰囲気で使ってる感じがしてよくないなと思ったので、色々調べた。

今回まとめたのは、MySQLのInnoDBのトランザクションについての一部。
この記事でまとめる対象の事柄は、既に詳しく素晴らしい記事が存在するが、それなりに難しい内容なのでそれらを読みながら理解の補助的に、またメモ的にまとめたものが本記事。

また、「この概念を理解するには、事前に〇〇の概念の理解が必要で、それを理解するにはこの記事がわかりやすい」といった形で自分が後から振り返られるようにもまとめた。

この記事はあくまでも理解の補助と、どの概要を学んだ方が良いかということをまとめた記事なので、詳細な理解は紹介している記事を読んだ方が良い。

# トランザクションの分離レベル
トランザクションには[ACID属性](https://ja.wikipedia.org/wiki/ACID_(%E3%82%B3%E3%83%B3%E3%83%94%E3%83%A5%E3%83%BC%E3%82%BF%E7%A7%91%E5%AD%A6))というものがあって、その中のI(Isolation)が分離性と呼ばれるもの。

トランザクションの分離にはレベルが存在する。
詳細は以下の記事等を参照。

大事なのは、どの分離レベルにおいて、どんな不都合な読み込みが行われるかを認識することだと思う。

[トランザクション分離レベル - Wikipedia](https://ja.wikipedia.org/wiki/%E3%83%88%E3%83%A9%E3%83%B3%E3%82%B6%E3%82%AF%E3%82%B7%E3%83%A7%E3%83%B3%E5%88%86%E9%9B%A2%E3%83%AC%E3%83%99%E3%83%AB)

[[RDBMS][SQL]トランザクション分離レベルについて極力分かりやすく解説 - Qiita](https://qiita.com/PruneMazui/items/4135fcf7621869726b4b)

# MySQLのInnoDBの Repeatable Read
MySQLのデフォルトのエンジンであるInnoDBのデフォルトのトランザクションの分離レベルはRepeatable Readである。

一般的に、Repeatable Readは、ダーティリードとファジーリードは防ぐことができるけれども、ファントムリードは防ぐことができない。

ファントムリードは、以下のようだとある。
> 別のトランザクションで挿入されたデータが見えることにより、一貫性がなくなる現象。

引用元: [[RDBMS][SQL]トランザクション分離レベルについて極力分かりやすく解説 - Qiita](https://qiita.com/PruneMazui/items/4135fcf7621869726b4b)

一般的なRepeatable Readの場合のシミュレーションすると以下のような感じ(実際にやったわけではない)

事前に以下のような操作を行なっているとする。

```sql
CREATE TABLE programming_langs (
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    name VARCHAR(20)
);

show tables;

+-------------------+
| Tables_in_sample  |
+-------------------+
| programming_langs |
+-------------------+

INSERT INTO programming_langs (name) VALUES ('Go');
```

Repeatable Readの場合のトランザクションのシミュレーション

![phantom.png](https://qiita-image-store.s3.amazonaws.com/0/145611/7fc400f8-191d-05f6-57f9-5d4e00d92bf2.png)

上記のように、左側のトランザクションがInsertしたデータを右側のトランザクションが読み込めてしまう。(ファントムリード)

しかし、MySQLの InnoDB は、 Repeatable Read であっても、別のトランザクションで挿入されたデータをうまいことやりくりしており、一貫性があるように見える。

実際にやってみた。

事前に以下のような操作を行なっているとする。

```sql
CREATE TABLE programming_langs (
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    name VARCHAR(20)
);

show tables;

+-------------------+
| Tables_in_sample  |
+-------------------+
| programming_langs |
+-------------------+

INSERT INTO programming_langs (name) VALUES ('Go');
```

Repeatable ReadのMySQL InnoDBでの操作を図にしたもの
(実際にMySQLでやってみたが、見易さを考慮して図にしている)

![mvcc.png](https://qiita-image-store.s3.amazonaws.com/0/145611/3509cc71-e9ea-e3d6-4d74-52857ca1f71f.png)

ファントムリードが生じていないように見える。


これはなぜだろうか。実は、[MVCC(MultiVersion Concurrency Control)](https://ja.wikipedia.org/wiki/MultiVersion_Concurrency_Control)という仕組みを用いているからである。

## MVCC
[MySQLのドキュメント](https://dev.mysql.com/doc/refman/8.0/en/glossary.html#glos_mvcc)によれば

> Acronym for “multiversion concurrency control”. This technique lets InnoDB transactions with certain isolation levels perform consistent read operations; that is, to query rows that are being updated by other transactions, and see the values from before those updates occurred. This is a powerful technique to increase concurrency, by allowing queries to proceed without waiting due to locks held by the other transactions.

> This technique is not universal in the database world. Some other database products, and some other MySQL storage engines, do not support it.

引用元: [MySQLのドキュメント](https://dev.mysql.com/doc/refman/8.0/en/glossary.html#glos_mvcc)

上記を意訳しつつ、まとめると以下のような感じ。

- MVCCを使用すると、特定のトランザクションの分離レベルのInnoDBで、一貫性のある読み取り操作が可能に
    - 他のトランザクションが更新した行を照会する際に、更新前の状態を見ることができる
- 他のトランザクションのロックを待たずに、クエリを進めることができるのでMVCCだと、同時実行性を高めることができる

> 他のトランザクションが更新した行を照会する際に、更新前の状態を見ることができる

これがまさにMySQLのInnoDBではRepeatable Readであっても別のトランザクションで挿入されたデータをうまいことやりくりして、一貫性があるように見える(ファントムリードが起きない)ことに寄与している。

これはファントムリードでは、トランザクション内で挿入したデータを別のトランザクション内で見えてしまうが、それを回避することができるということだ。

MVCCの仕組みの説明は、以下がわかりやすいのでこちらを参照。
[MySQLのMVCC - Qiita](https://qiita.com/nkriskeeic/items/24b7714b749d38bba87b)

[公式](https://dev.mysql.com/doc/refman/8.0/en/innodb-multi-versioning.html)でも詳しく仕組みが記されている。

# MVCCが故の注意

行ロックをすることなく一貫性を出せるのと、複数のトランザクションが同一のテーブルを操作するとき、お互いのロックを待つことなく操作できるので、同時実行性が上がる。

しかし、行ロックをしないが故に生じ得るLost Updateという問題も存在する。
以下に詳しく記述されているので参照。
[漢(オトコ)のコンピュータ道: InnoDBのREPEATABLE READにおけるLocking Readについての注意点](http://nippondanji.blogspot.com/2013/12/innodbrepeatable-readlocking-read.html)

詳しくは上記の記事を参照した方が良いが、記事を読んだ際に頭を整理するためにこれまで書いてきた内容と合わせて、以下のようなメモをしながら読んだので載せておく。

- Lost Update
    - MVCCは、行ロックをとるわけではないので2つのトランザクションが重なったときには、後で更新した方の値が最終的に反映される
- Locking Read
    - Lost Updateを防ぐためには、Locking Readを使用する
        - Lost Update は行にロックをかけていないから生じるわけで、それならば行をロックすればええやんという感じ
            - 排他ロックと共有ロックを行える構文が存在する
                - 排他ロックと共有ロックは[この記事](https://qiita.com/mizzwithliam/items/31fb68217899bd0559e8)がわかりやすい
        - 逆説的な話だが、MVCCでRepeatable Readでもファントムリードを防げているのは、MVCCによって、DB_ROLL_PTRの値を見ているからで、Locking Readの場合はそうでなくファントムリードが生じる
            - > Locking ReadはREAD COMMITTED 
            - > UPDATEとDELETEやデフォルトでLocking Readの挙動
                - (引用元: [漢(オトコ)のコンピュータ道: InnoDBのREPEATABLE READにおけるLocking Readについての注意点](http://nippondanji.blogspot.com/2013/12/innodbrepeatable-readlocking-read.html))
     

> DB_ROLL_PTR ... そのレコードの過去の値を持つundo log recordへのポインタ

引用元: [MySQLのMVCC - Qiita](https://qiita.com/nkriskeeic/items/24b7714b749d38bba87b)


# ロックの仕組み
MySQLのロックにはレコードロック、ギャップロック、ネクストキーロックの３つがある。
以下の記事が非常にわかりやすいので、参照。

[MySQL(InnoDB)のネクストキーロックの仕組みと範囲を図解する - 備忘録の裏のチラシ](https://norikone.hatenablog.com/entry/2018/09/12/MySQL%28InnoDB%29%E3%81%AE%E3%83%8D%E3%82%AF%E3%82%B9%E3%83%88%E3%82%AD%E3%83%BC%E3%83%AD%E3%83%83%E3%82%AF%E3%81%AE%E4%BB%95%E7%B5%84%E3%81%BF%E3%81%A8%E7%AF%84%E5%9B%B2%E3%82%92%E5%9B%B3%E8%A7%A3)

ただ、「レコードロックっていうのは実際は、インデックスレコードのロックだよ」みたいな感じでインデックスの知識が必要。
ここで必要なインデックスの知識も、ただクエリの性能向上のための便利なものという理解ではなく内部の仕組みの理解が必要みたいだ。(特にクラスタインデックスとセカンダリインデックス)

なので、上記のロックの記事を読む前に以下の記事等でインデックスの仕組みについて学んでからの方が良さそう。

インデックスの記事
[漢(オトコ)のコンピュータ道: 知って得するInnoDBセカンダリインデックス活用術！](http://nippondanji.blogspot.com/2010/10/innodb.html)

[INDEX FULL SCANを狙う](https://sh2.hatenablog.jp/entries/2011/12/17)

[MySQL :: MySQL 5.6 リファレンスマニュアル :: 14.2.13.2 クラスタインデックスとセカンダリインデックス](https://dev.mysql.com/doc/refman/5.6/ja/innodb-index-types.html)

[MySQL with InnoDB のインデックスの基礎知識とありがちな間違い - クックパッド開発者ブログ](https://techlife.cookpad.com/entry/2017/04/18/092524)

[MySQL(InnoDB)のインデックスについての備忘録 - What is it, naokirin?](https://naokirin.hatenablog.com/entry/2015/02/07/193609)


# 参考にさせていただいたサイト

## ACID属性
[ACID (コンピュータ科学) - Wikipedia](https://ja.wikipedia.org/wiki/ACID_(%E3%82%B3%E3%83%B3%E3%83%94%E3%83%A5%E3%83%BC%E3%82%BF%E7%A7%91%E5%AD%A6))

## トランザクション分離レベル
[[RDBMS][SQL]トランザクション分離レベルについて極力分かりやすく解説 - Qiita](https://qiita.com/PruneMazui/items/4135fcf7621869726b4b)

[トランザクション分離レベル - Wikipedia](https://ja.wikipedia.org/wiki/%E3%83%88%E3%83%A9%E3%83%B3%E3%82%B6%E3%82%AF%E3%82%B7%E3%83%A7%E3%83%B3%E5%88%86%E9%9B%A2%E3%83%AC%E3%83%99%E3%83%AB)

## MVCC
[MySQLのMVCC - Qiita](https://qiita.com/nkriskeeic/items/24b7714b749d38bba87b)

[MySQLのドキュメント](https://dev.mysql.com/doc/refman/8.0/en/glossary.html#glos_mvcc)

ロック
[MySQLでSELECT FOR UPDATEと行ロックの挙動を検証してみた - JUST FOR FUN](https://taiga.hatenadiary.com/entry/2018/02/12/170109)

[MySQL - InnoDBのロック関連まとめ - Qiita](https://qiita.com/mizzwithliam/items/31fb68217899bd0559e8)

[世界の何処かで MySQL（InnoDB）の REPEATABLE READ に嵌る人を1人でも減らすために - KAYAC engineers' blog](https://techblog.kayac.com/repeatable_read.html)

[doc/innodb.md at master · ichirin2501/doc](https://github.com/ichirin2501/doc/blob/master/innodb.md)

[漢(オトコ)のコンピュータ道: InnoDBのREPEATABLE READにおけるLocking Readについての注意点](http://nippondanji.blogspot.com/2013/12/innodbrepeatable-readlocking-read.html)

[アプリケーションエンジニアが知っておくべきMySQLのロック - Qiita](https://qiita.com/tikamoto/items/f867050ff77d06a94215)

[MySQL(InnoDB)のネクストキーロックの仕組みと範囲を図解する - 備忘録の裏のチラシ](https://norikone.hatenablog.com/entry/2018/09/12/MySQL%28InnoDB%29%E3%81%AE%E3%83%8D%E3%82%AF%E3%82%B9%E3%83%88%E3%82%AD%E3%83%BC%E3%83%AD%E3%83%83%E3%82%AF%E3%81%AE%E4%BB%95%E7%B5%84%E3%81%BF%E3%81%A8%E7%AF%84%E5%9B%B2%E3%82%92%E5%9B%B3%E8%A7%A3)

インデックス
[漢(オトコ)のコンピュータ道: 知って得するInnoDBセカンダリインデックス活用術！](http://nippondanji.blogspot.com/2010/10/innodb.html)

[INDEX FULL SCANを狙う](https://sh2.hatenablog.jp/entries/2011/12/17)

[MySQL :: MySQL 5.6 リファレンスマニュアル :: 14.2.13.2 クラスタインデックスとセカンダリインデックス](https://dev.mysql.com/doc/refman/5.6/ja/innodb-index-types.html)

[MySQL with InnoDB のインデックスの基礎知識とありがちな間違い - クックパッド開発者ブログ](https://techlife.cookpad.com/entry/2017/04/18/092524)

[MySQL(InnoDB)のインデックスについての備忘録 - What is it, naokirin?](https://naokirin.hatenablog.com/entry/2015/02/07/193609)
