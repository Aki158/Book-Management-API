# School-API

## 概要

教師IDをもとに生徒やクラスの情報を取得できるWebAPI

## 実装方針

このプロジェクトの目的は、シンプルで拡張性のあるWeb APIを提供することです。Goを使用したWebサーバーとMySQLデータベースをDockerコンテナ上で動作させることで、ローカル環境でWebAPIを利用できます。

## アーキテクチャ

このプロジェクトは、以下のようなアーキテクチャで設計されています。

- Webサーバー: Goで実装されており、RESTful APIを提供します。ユーザーからのリクエストを受け取り、MySQLデータベースと連携してデータの管理を行います。
- データベースサーバー: MySQLを使用しており、データベーススキーマはシンプルで拡張性を考慮した設計になっています。
- Docker: 各サービスはDockerコンテナとして実行され、コンテナ間の通信はDockerネットワークを介して行います。

## 技術スタック

| カテゴリ | 使用技術 |
| ------- | ------- |
| 言語 | Go |
| データベース | MySQL |
| インフラ | Docker |

## ER図

ER図は、以下3つのテーブルで構成しています。

- teachers
- students
- classes

実際のER図は、[er.md](https://github.com/Aki158/School-API/blob/main/design/er.md)を確認してください。

## セットアップ手順

ターミナル上で以下コマンドを実施することによりローカル環境で動作確認ができます。

1. リポジトリをクローンする
```
git clone https://github.com/Aki158/School-API.git
```

2. クローンしたリポジトリへ移動する
```
cd School-API
```

3. Dockerコンテナを起動する
```
docker compose up
```

### 使用方法

### エンドポイント

```bash
GET /students
```

### リクエスト

| パラメーター名 | 型 | 必須/任意 | デフォルト値 | 概要 |
| ------- | ------- | ------- | ------- | ------- |
| `facilitator_id` | int | 必須 |  | 教師ID |
| `page` | int | 任意 |  | ページネーションのページ数 |
| `limit` | int | 任意 | 5 | ページネーションの1ページあたりの表示数 |
| `sort` | string | 任意 | id | ソートのキーとしてレスポンスにある以下のフィールドを指定できる<br>・id: 生徒ID<br>・name: 生徒名<br>・loginId: ログインID |
| `order` | string | 任意 | asc | ソートの昇順=asc, 降順=descを指定する |
| `{key}_like` | string | 任意 |  | 指定したレスポンスのフィールドによる部分一致検索をします。<br>{key}にはレスポンスにある以下のフィールドを指定できる。<br>・name: 生徒名<br>・loginId: ログインID |

リクエスト例

```bash
curl -k -i 'https://127.0.0.1:48080/students?facilitator_id=1'
```

### レスポンス

| フィールド | 型 | 概要 |
| ------- | ------- | ------- |
| `$.students` | array |  |
| `$.students[*].id` | int | 生徒ID |
| `$.students[*].name` | string | 生徒名 |
| `$.students[*].loginId` | string | ログインID |
| `$.students[*].classroom` | array |  |
| `$.students[*].classroom[*].id` | int | クラスID |
| `$.students[*].classroom[*].name` | string | クラス名 |
| `$.totalCount` | int | リクエストの条件に該当する件数 |

レスポンス例

```json
{
    "students": [
        {
            "id": 1,
            "name": "佐藤",
            "loginId": "foo123",
            "classroom": {
                "id": 1,
                "name": "クラスA"
            }
        }
    ],
    "totalCount": 1
}
```

### エラー

以下の2種類に該当する場合にエラーを返します

| 状況 | レスポンスコード |
| ------- | ------- |
| リクエストに該当する生徒が存在しない | 404 Not Found |
| リクエストに問題がある | 400 Bad Request |

### 補足説明

- Dockerコンテナ起動直後は、MySQLコンテナの起動が完了していないことにより、リクエストに失敗することがあります。その場合は、時間をおいてから再度実施してください。
- このWebAPIは、最低限の手順でローカル環境での動作確認ができるようにしています。通常は、環境変数で管理するべき情報（ポート番号やDBへの接続情報）については、ハードコーディングしています。
