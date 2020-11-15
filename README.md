# Power Phrase

心に残る言葉を共有するサイトです。一覧画面でYouTube動画を再生することもできます。

<img width="800px" alt="トップページ" src="https://user-images.githubusercontent.com/61341861/99197759-c9154b80-27d7-11eb-8df9-e8270f7ae053.png">

# URL

https://www.power-phrase.com/

ヘッダーにある「動作確認用ログイン」をクリックすると、動作確認用ユーザーとしてログインできます。

# 使用技術
## バックエンド

- Go 1.15.2
    - Echo v3.3.10(RESTフレームワーク)
    - GORM v1.9.16(ORM)
    - validator.v9 v9.31.0
    - jwt-go v3.2.0(JWT認証)
    - go-imageupload v0.0.0-20160503070439-09d2b92fa05e
    - godotenv v1.3.0(環境変数ファイル読み込み)
    - air v1.21.2(ホットリロード)
    - testify v1.4.0(自動テスト)
- AWS(インフラ構成図参照)
    - EC2(Amazon Linux 2)
    - RDS(MySQL 8)
    - CloudFront
    - S3
    - Route53
    - ALB
    - ACM

## フロントエンド
- JavaScript
- Vue.js
    - Vue CLI 4.5.8
    - Vue Router 3.4.9
- Vuetify 2.3.17
- Axios 0.19.2
- npm 6.14.6

## 共通
- macOS Catalina 10.15.7
- Git 2.24.3
- SourceTree 4.0.2
- Docker Desktop for Mac 2.5.0.0(開発環境)
    - Docker 19.03.13
    - Docker Compose 1.27.4

# インフラ構成図

<img width="500" alt="インフラ構成図" src="https://user-images.githubusercontent.com/61341861/99197827-15f92200-27d8-11eb-8adc-c78c3756260c.png">

# 機能一覧

- 投稿一覧機能
- ページネーション機能
- YouTube動画再生機能
- 投稿検索機能
- ユーザー登録機能(プロフィール画像アップロード含む)
- ログイン機能(JWT認証)
- 動作確認用ログイン機能
- ログアウト機能
- 投稿登録機能(ログイン後のみ可能)
- 投稿詳細機能(ログイン後のみ可能)
- 投稿編集機能(ログイン後、自分が登録したものについてのみ可能)
- 投稿削除機能(ログイン後自分が登録したものについてのみ可能)
- ユーザー詳細表示機能(該当ユーザーによる投稿一覧含む)
- ユーザー更新機能
