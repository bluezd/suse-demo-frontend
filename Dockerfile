FROM golang:1.15-alpine as build-stage

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk --no-cache add build-base git mercurial gcc

WORKDIR /app

ADD go.mod go.sum logHandler.go main.go /app/

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod download \
    && go build

FROM alpine:latest

WORKDIR /app

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache ca-certificates

COPY --from=build-stage /app/suse-demo-frontend /app/
COPY templates /app/templates

EXPOSE 8000

CMD ["./suse-demo-frontend"]
