FROM golang:1.23

# 作業ディレクトリを設定
WORKDIR /app

# bashとopensslをインストール
RUN apt-get update && apt-get install -y bash openssl

# アプリケーションのソースコードをコピー
COPY main.go /app
COPY go.mod  /app
# COPY go.sum  /app

# 依存関係をダウンロードして go.sum を更新
# RUN go mod tidy
# RUN go mod download

# 自己署名証明書を生成する
RUN bash -c 'openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj "/CN=localhost"'

# 依存関係のチェックを含むビルドステップを実施
RUN go build -o server

# コンテナ起動時にmain.goを実行
CMD ["/app/server"]