FROM golang:1.17-alpine as builder
WORKDIR /build

ARG GOARCH=amd64

COPY go.mod .
COPY go.sum .
COPY . .

RUN GOOS=linux GOARCH=$GOARCH GO111MODULE=on go build -v -o ./build/run cmd/using-gin/main.go

FROM alpine
WORKDIR /app
ENV USING_GIN_ENV=production

COPY --from=builder /build/.env .
COPY --from=builder /build/build/run .

ENTRYPOINT ["./run"]

