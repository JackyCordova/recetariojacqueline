FROM golang:1.24.6 AS builder

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -o /workspace/recipe-server ./recipe/cmd/grpcmain

FROM gcr.io/distroless/base-debian12

COPY --from=builder /workspace/recipe-server /recipe-server

EXPOSE 9093

ENTRYPOINT ["/recipe-server"]