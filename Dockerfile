FROM golang:1.23.0 as builder

WORKDIR /app

ENV GOPROXY=https://mirrors.tencent.com/go/

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloud-bill main.go

FROM alpine:latest

WORKDIR /app

RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update --no-cache \
    && apk --no-cache add ca-certificates tzdata 

ENV TZ=Asia/Shanghai

COPY --from=builder /app/cloud-bill .
COPY --from=builder /app/config-tmpl.yaml ./config.yaml

CMD ["./cloud-bill"]
