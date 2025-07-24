# Go Todo App - Alpine版 Dockerfile（学習最適化版）

# ===== Build Stage =====
FROM golang:1.21-alpine AS builder

# セキュリティ＆開発に必要なパッケージを追加
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata

# 作業ディレクトリを設定
WORKDIR /app

# Go modulesファイルを先にコピー（レイヤーキャッシュ最適化）
COPY go.mod go.sum ./

# 依存関係をダウンロード＆検証
RUN go mod download && go mod verify

# ソースコードをコピー
COPY . .

# 静的バイナリをビルド
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s' \
    -o main .

# ===== Runtime Stage =====
FROM alpine:latest

# ランタイムに必要なパッケージを追加
RUN apk add --no-cache \
    ca-certificates \
    wget \
    curl \
    tzdata

# セキュリティ：非rootユーザーを作成
RUN addgroup -S appuser && adduser -S appuser -G appuser

# 作業ディレクトリを設定
WORKDIR /app

# ビルドステージから実行ファイルをコピー
COPY --from=builder /app/main .

# ファイルの所有者を変更
RUN chown -R appuser:appuser /app

# 非rootユーザーに切り替え
USER appuser

# アプリケーションが使用するポートを公開
EXPOSE 8080

# ヘルスチェックを設定
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# アプリケーションを起動
CMD ["./main"]