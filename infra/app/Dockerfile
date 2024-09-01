FROM golang:1.23

# 作業ディレクトリを設定
WORKDIR /app

# bashをインストール
RUN apt-get update && apt-get install -y bash

# アプリケーションのソースコードをコピー
COPY main.go /app
COPY go.mod  /app
# COPY go.sum  /app

# 依存関係をダウンロードして go.sum を更新
# RUN go mod tidy
# RUN go mod download

# 依存関係のチェックを含むビルドステップを実施
RUN go build -o server

# コンテナ起動時にmain.goを実行
CMD ["/app/server"]