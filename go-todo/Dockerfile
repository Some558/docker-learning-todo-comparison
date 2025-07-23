# Go Todo App Dockerfile - 修正版

# ===== Build Stage =====
FROM golang:1.21-alpine AS builder

# セキュリティアップデート + 必要ツール
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Go modules（キャッシュ効率化）
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
FROM alpine:latest

# セキュリティ: 証明書追加
RUN apk --no-cache add ca-certificates wget

# セキュリティ: 非rootユーザー作成
RUN addgroup -S appuser && adduser -S appuser -G appuser

# 作業ディレクトリ設定
WORKDIR /app

# ビルドステージから実行ファイルをコピー
COPY --from=builder /app/main .

# 所有者を変更（セキュリティ強化）
RUN chown -R appuser:appuser /app

# 非rootユーザーで実行
USER appuser

# ポート公開
EXPOSE 8080

# ヘルスチェック設定
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 実行
ENTRYPOINT ["./main"]