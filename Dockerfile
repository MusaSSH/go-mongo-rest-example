FROM golang:1.19.0-alpine3.16

ENV MONGOURI="mongodb://"

WORKDIR /usr/src/app

COPY . .

RUN go mod download && go mod verify && \
    go build -v -o /usr/local/bin/app

CMD [ "app" ]