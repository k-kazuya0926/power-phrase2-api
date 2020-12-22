# Power Phrase

心に残る言葉を共有するサイトです。

以前精神的に辛かった時期に言葉の大切さに気づいた経験から、良い言葉を共有することにより、
利用者がより前向きな気持ちになれると良いという思いから作成しました。

<img width="800px" alt="トップページ" src="https://user-images.githubusercontent.com/61341861/102834081-5d597a80-4436-11eb-80e0-e043c30c4f93.png">

# URL

https://www.power-phrase.com/

ヘッダーにある「動作確認用ログイン」をクリックすると、動作確認用ユーザーとしてログインできます。

# 特徴

## バックエンド

- DDD(オニオンアーキテクチャ)
- モックやFour Phase Testを取り入れた自動テスト
- AWS ECSによるコンテナデプロイ
- Terraformによるインフラのコード化(VPC、RDS)
- CircleCIによるCI

## フロントエンド

- Vue.jsによるシングルページアプリケーション
- レスポンシブデザイン
- AWS CloudFront、S3による配信
- CircleCIによるCD

# 使用技術
## バックエンド

- Go 1.15.2
    - DDD(オニオンアーキテクチャ)
    - Echo v3.3.10(RESTフレームワーク)
    - GORM v1.9.16(ORM)
    - validator.v9 v9.31.0
    - jwt-go v3.2.0(JWT認証)
    - go-imageupload v0.0.0-20160503070439-09d2b92fa05e
    - godotenv v1.3.0(環境変数ファイル読み込み)
    - air v1.21.2(ホットリロード)
    - testify v1.4.0(自動テスト)
- MySQL 8.0.21
- AWS(下記インフラ構成図参照)
    - ECS
    - ECR
    - EC2(Amazon Linux 2)
    - RDS(MySQL 8)
    - Route53
    - ALB
    - ACM
- Terraform v0.13.5
- CircleCI

## フロントエンド(https://github.com/k-kazuya0926/power-phrase2-front)
- JavaScript
- Vue.js 2.6.12
    - Vue CLI 4.5.9
    - Vue Router 3.4.9
    - Vuex 3.5.1
    - vuex-persistedstate 4.0.0-beta.1(Vuexの永続化)
    - Vuetify 2.3.17(UIフレームワーク)
    - Axios 0.19.2(Ajax)
    - VeeValidate 3.4.5
    - moment 2.29.1(日付操作)
- Node.js 12.18.4
- npm 6.14.6
- AWS(下記インフラ構成図参照)
    - CloudFront
    - S3
    - Route53
    - ACM
- CircleCI

## 共通
- macOS Catalina 10.15.7
- Git 2.24.3
- SourceTree 4.0.2
- Docker Desktop for Mac 2.5.0.0(開発環境)
    - Docker 19.03.13
    - Docker Compose 1.27.4

# インフラ構成図

<img width="640px" alt="インフラ構成図" src="https://user-images.githubusercontent.com/61341861/102556459-494b0b80-410c-11eb-8877-bd032cfc6110.jpg">

# 機能一覧

- 投稿一覧機能
- ページネーション機能
- YouTube動画再生機能
- 投稿検索機能(タイトル、発言者、詳細のいずれかがキーワードを含むという条件での検索)
- ユーザー登録機能(プロフィール画像アップロード含む)
- ログイン機能(JWT認証)
- 動作確認用ログイン機能
- ログアウト機能
- ユーザー詳細表示機能(該当ユーザーによる投稿一覧含む)
- ユーザー更新機能
- ユーザー削除機能
- 投稿登録機能(ログイン後のみ可能)
- 投稿詳細機能(ログイン後のみ可能)
- 投稿更新機能(ログイン後、自分が登録したものについてのみ可能)
- 投稿削除機能(ログイン後自分が登録したものについてのみ可能)
- コメント登録機能
- コメント一覧機能
- コメント削除機能(ログイン後、自分が登録したものについてのみ可能)
