FROM golang:latest

ENV TODO_PORT=7540
ENV TODO_DBFILE=/scheduler.db
ENV TODO_PASSWORD=primite-plz

EXPOSE ${TODO_PORT}

WORKDIR /app

COPY . /app

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main

CMD ["/main"]

