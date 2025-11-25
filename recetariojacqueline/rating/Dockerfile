FROM golang:1.24.6 AS builder

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -o /workspace/rating-server ./rating/cmd/grpcmain

FROM gcr.io/distroless/base-debian12

COPY --from=builder /workspace/rating-server /rating-server

EXPOSE 9092

ENTRYPOINT ["/rating-server"]