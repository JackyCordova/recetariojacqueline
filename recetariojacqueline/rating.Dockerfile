# imagen de go que vamos a usar
FROM golang:1.24

# en que carpeta se hará todo
WORKDIR /usr/src/app

# copiamos los módulos de go y descargamos las librerías necesarias
COPY go.mod go.sum ./
RUN go mod download

# copiamos el código fuente
COPY ./pkg ./pkg
COPY ./rating ./rating

# cual es el comando que tiene que correr cuando se inicie
CMD ["go", "run", "rating/cmd/main/main.go"]
