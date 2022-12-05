FROM golang:1.19.0-alpine3.15

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o app
RUN chmod +x ./app
RUN mkdir -m 0777 /storage

EXPOSE 8080

ENTRYPOINT [ "./app" ]
