FROM golang:1.22.4 as builder

WORKDIR /app

COPY . .

RUN env CGO_ENABLED=0 GOOS=linux go build -o gitlab-metadata

FROM alpine

RUN apk add --no-cache git

WORKDIR /app

COPY --from=builder /app/gitlab-metadata /app/gitlab-metadata

CMD ["/app/gitlab-metadata"]

