FROM golang:1.24.6 AS builder

WORKDIR /workspace

# Copiamos go.mod y go.sum del proyecto raíz
COPY go.mod go.sum ./
RUN go mod download

# Copiamos TODO el proyecto (importantísimo)
COPY . .

# Compilar binario de metadata
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -o /workspace/metadata-server ./metadata/cmd/grpcmain

# Imagen final
FROM gcr.io/distroless/base-debian12

COPY --from=builder /workspace/metadata-server /metadata-server

EXPOSE 9091

ENTRYPOINT ["/metadata-server"]