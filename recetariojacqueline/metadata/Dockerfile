#dockerfile - nos ayuda a construir la imagen de docker
#imagen de golang que vamos a utilizar
FROM golang:1.24 AS build 
#el directiorio de trabajo del contenedor es /app
WORKDIR /app
# copiamos go.mod y go.sum al contenedor, para definir y descargar las dependencias
COPY go.mod go.sum ./
RUN go mod download
COPY . .
#compilamos el binario de metadata
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/metadata ./metadata/cmd/main

FROM gcr.io/distroless/base-debian12:nonroot
#variables de entorno
ENV PORT=8081 CONSUL_ADDR=consul:8500 SERVICE_HOST=metadata
COPY --from=build /out/metadata /metadata
#en que puerto está escuchando el servicio
EXPOSE 8081
USER nonroot:nonroot
#comando principal del contenedor
ENTRYPOINT ["/metadata"]
