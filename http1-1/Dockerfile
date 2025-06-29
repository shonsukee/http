FROM golang:1.23

WORKDIR /app

# 必要なファイルをコピー
COPY . /app/

# モジュールの初期化
RUN go mod init & go mod tidy

EXPOSE 8080

CMD ["go", "run", "."]