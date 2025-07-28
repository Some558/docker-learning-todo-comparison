# Dockerfile.level2 - Distroless版
# 目標: 55.4MB → 25MB（劇的軽量化）

# ===== Build Stage =====
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Go modulesを先にコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download && go mod verify

# ソースコードをコピー
COPY . .

# 完全静的バイナリビルド（Distroless対応）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -trimpath \
    -tags 'netgo osusergo' \
    -o main .

# バイナリサイズ確認
RUN ls -lh main

# ===== Distroless Runtime Stage =====
FROM gcr.io/distroless/static:nonroot

# CA証明書をコピー（HTTPS通信用）
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# アプリケーションバイナリをコピー
COPY --from=builder /app/main /app/main

EXPOSE 8080

# nonrootユーザーで実行（セキュア）
USER 65532:65532

ENTRYPOINT ["/app/main"]