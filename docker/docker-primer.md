# Dockerの最低限の知識と流れをサクッとまとめる

## 本記事の目的
大きく分けて2つ。

* Dockerを初めて行くのにあたって最低限必要な知識をサクッと理解する
仕組みや詳細については、参照サイトや書籍を提示し、そこを見ればわかるようにする。

* Dockerを初めて行くのにあたって最低限理解が必要な流れをサクッと理解する
大まかな流れを理解して、詳細や実際のコマンドは、チートシート等を見ればわかるようにする。

## 最低限の知識
### 仮想化の種類
Dockerの仮想化はコンテナ型。
[「Docker」を全く知らない人のために「Docker」の魅力を伝えるための「Docker」入門 - Qiita](https://qiita.com/bremen/items/4604f530fe25786240db#docker%E3%81%AE%E4%BB%95%E7%B5%84%E3%81%BF)で詳しく解説されているので、参照されたし。
この中で特に大事なのは、以下の部分。

> コンテナはホストOSから見ると単一のプロセスとして扱われ、カーネル部分をホストOSと共有するため、リソース使用量が非常に少ない

引用元 : [「Docker」を全く知らない人のために「Docker」の魅力を伝えるための「Docker」入門 - Qiita](https://qiita.com/bremen/items/4604f530fe25786240db#docker%E3%81%AE%E4%BB%95%E7%B5%84%E3%81%BF)


これによって他の仮想化よりも、軽量に動作させることができる。
ただ、「ホストカーネルを共有する」と言うところがマイナスに働く面もあり、最近それを解消するためのgVisorと言うのが発表されて話題になっている。
[20分でわかるgVisor入門](https://www.slideshare.net/uzy_exe/201805gvisorintroduciton)

### コンテナ
> コンテナは、ホストOS上に作成された論理的な区画

引用元 : WINGSプロジェクト阿佐志保　(2018/4/11)『プログラマのためのDocker教科書 第2版 インフラの基礎知識&コードによる環境構築の自動化』 翔泳社

ホストOS上に起動する仮想的なサーバーのようなもので、このコンテナ上でインストールした諸々のライブリ等を動作させる。

### イメージ
上記のコンテナを起動し、実行する際に必要なライブラリ等のファイルシステムで設定がまとめられている。
Read-Onlyである。
[Docker Hub](https://hub.docker.com/) (後述)に公式や誰かが作成したイメージが公開されている。

そのDockerのイメージを元に自分で必要なイメージをさらに追加して自分用のイメージを作成していく感じ。
イメージはレイヤ構成になっていて、必要なイメージを重ねていく。
公式では、[nginx](https://hub.docker.com/_/nginx/)や[centos](https://hub.docker.com/_/centos/)と言ったミドルウェアやOSのイメージが公開されている。

例えば、OSとしてubuntuを使用し、その上にGoとMySQLを使用したアプリケーションを作成する場合、ubuntu、Go、MySQL等のそれぞれのイメージをDocker Hubから取得し、そのイメージを使用して実行環境のイメージを作成すると言った具合で使用して行く。

### Dockerfile
イメージの設定ファイル(構成情報)。
これを書くことによってインフラの設定をコードで管理することができるようになる。
Docker Hub等からインストールした公式のイメージ等を組み合わせたりして、自分でイメージを作成することができるが、Dockerfileをビルドすることによって、環境毎に面倒臭い諸々の設定等を毎回行うことなく、他の環境でも同じイメージを作成し、コンテナを起動することができる。つまり、1回Dockerfileを書いてしまえば、面倒臭いインフラの環境構築を毎回する必要がなくなるのだ。

### Docker Hub
Dockerイメージを管理するリポジトリ。Git Hubみたいなもの。

## 最低限の流れ
イメージ取得、イメージの作成、コンテナ起動、イメージをDocker Hubにpushする等の大まかな流れを記述する。
ここで記載するコマンドは基本的なものに留める。
オプション等の詳細は、他の記事参照。

### 既に存在するものをそのまま使用する
1 Docker Hubからコンテナの元となるイメージを取得する

```
docker image pull imageName:tagName
```

2 取得したイメージを元にコンテナを起動
→　pullを飛ばしてdocker runだけでも可能、ローカルに存在しない場合はDockerが自動でpullしてくれる

```
docker run imageName
```

### 新しくイメージを自分で作成する場合(コンテナから作成する場合)
1 Docker Hubからコンテナの元となるイメージを取得する

```
docker image pull imageName:tagName
```

2 Dockerイメージを元にコンテナを起動し、諸々の設定を手動で行い、そのコンテナを元にイメージを作成する

```
docker commit containerName imageName
```

参考 : [Dockerでcommitしてみる - Qiita](https://qiita.com/mats116/items/712575dc50513dfdf0a2)

### 新しく自分で作成する場合(Dockerfileから作成する場合)
1 Docker Hubからコンテナの元となるイメージを取得する

```
docker image pull imageName:tagName
```

2 Dockerイメージや、諸々の設定をDockerfile(設定ファイル)に記述する  
[【入門】Dockerfileの基本的な書き方 | レコチョクのエンジニアブログ](https://techblog.recochoku.jp/1022) 等が参考になる

3 Dockerfileを元にイメージをビルドする(docker build)

```
docker build -t imageName:tagName path/to/Dockerfile
```

4 作成したイメージからコンテナを起動する(docker run)

```
docker run imageName
```

5 作成したイメージをDocker Hubにpushする  
コードをgit hubにあげるようなもの  
[docker push](https://qiita.com/suin/items/20d735823e158196983e)

## 今後の学習の流れ
### 仕組みの詳細が知りたい
WINGSプロジェクト阿佐志保　(2018/4/11)『プログラマのためのDocker教科書 第2版 インフラの基礎知識&コードによる環境構築の自動化』 翔泳社

[ゼロからはじめる Dockerによるアプリケーション実行環境構築 | Udemy](https://www.udemy.com/docker-k/)

### 実際のコマンドが知りたい
#### 網羅的に学びたい
WINGSプロジェクト阿佐志保　(2018/4/11)『プログラマのためのDocker教科書 第2版 インフラの基礎知識&コードによる環境構築の自動化』 翔泳社

[ゼロからはじめる Dockerによるアプリケーション実行環境構築 | Udemy](https://www.udemy.com/docker-k/)

#### リファレンスとして参照したい
[docker コマンド チートシート - Qiita](https://qiita.com/voluntas/items/68c1fd04dd3d507d4083)

[Dockerコマンドメモ - Qiita](https://qiita.com/curseoff/items/a9e64ad01d673abb6866)

### Dockerfileを実際に書きたい
WINGSプロジェクト阿佐志保　(2018/4/11)『プログラマのためのDocker教科書 第2版 インフラの基礎知識&コードによる環境構築の自動化』 翔泳社

[Dockerfileを書いてみる - Qiita](https://qiita.com/nl0_blu/items/1de829288db2670276e8)

[【入門】Dockerfileの基本的な書き方 | レコチョクのエンジニアブログ](https://techblog.recochoku.jp/1022)

## 参考文献
WINGSプロジェクト阿佐志保　(2018/4/11)『プログラマのためのDocker教科書 第2版 インフラの基礎知識&コードによる環境構築の自動化』 翔泳社

[ゼロからはじめる Dockerによるアプリケーション実行環境構築 | Udemy](https://www.udemy.com/docker-k/)

## 参考にさせていただいたサイト
[Dockerfileを書いてみる - Qiita](https://qiita.com/nl0_blu/items/1de829288db2670276e8)

[【入門】Dockerfileの基本的な書き方 | レコチョクのエンジニアブログ](https://techblog.recochoku.jp/1022)

[docker コマンド チートシート - Qiita](https://qiita.com/voluntas/items/68c1fd04dd3d507d4083)

[Dockerコマンドメモ - Qiita](https://qiita.com/curseoff/items/a9e64ad01d673abb6866)

[「Docker」を全く知らない人のために「Docker」の魅力を伝えるための「Docker」入門 - Qiita](https://qiita.com/bremen/items/4604f530fe25786240db#docker%E3%81%AE%E4%BB%95%E7%B5%84%E3%81%BF)

[Dockerイメージの理解とコンテナのライフサイクル](https://www.slideshare.net/zembutsu/docker-images-containers-and-lifecycle)

[Dockerでcommitしてみる - Qiita](https://qiita.com/mats116/items/712575dc50513dfdf0a2)
