FROM golang:alpine

WORKDIR /app/
ADD . /app/

RUN apk update && apk upgrade \
    && apk add --no-cache git

RUN go mod download

ARG SERVICE_NAME

RUN go build -o bin/app ./cmd/$SERVICE_NAME/main.go

ENTRYPOINT ["bin/app"]