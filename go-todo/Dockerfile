FROM golang:1.21-alpine AS builder

# セキュリティアップデート + 必要ツール
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Go modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# ソースコード
COPY . .

# 最適化ビルド
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main .

# ===== Runtime Stage =====
FROM scratch

# CA証明書とタイムゾーン情報をコピー
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# アプリケーション
COPY --from=builder /app/main /main

# ポート公開
EXPOSE 8080

# ヘルスチェック用のエンドポイント
# (scratchイメージではwgetが使えないので、アプリ側で実装)

# 実行
ENTRYPOINT ["/main"]
