FROM golang:1.16-alpine as builder
WORKDIR /build

COPY go.mod .
COPY go.sum .
COPY . .

RUN GOOS=linux GOARCH=amd64 GO111MODULE=on go build -v -o ./build/run cmd/using-gin/main.go

FROM alpine
WORKDIR /app
ENV USING_GIN_ENV=production

COPY --from=builder /build/.env .
COPY --from=builder /build/build/run .

ENTRYPOINT ["./run"]

