FROM golang:1.13.5
WORKDIR /go/src/github.com/augmentable-dev/deno-run-http
COPY main.go main.go
COPY go.mod go.mod
# COPY go.sum go.sum
RUN go build -o server main.go

FROM hayd/alpine-deno:1.0.0
RUN apk update && apk upgrade && \
    apk add --no-cache bash
COPY --from=0 /go/src/github.com/augmentable-dev/deno-run-http/server server

EXPOSE 8000

ENTRYPOINT [ "./server" ]
