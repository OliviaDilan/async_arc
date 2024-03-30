FROM golang:alpine

RUN apk update && apk upgrade \
    && apk add --no-cache git

ARG SERVICE_NAME

WORKDIR /services/
ADD ./${SERVICE_NAME}/ /services/app/
ADD ./pkg/ /services/pkg/

WORKDIR /services/app/

RUN go mod download

RUN go build -o bin/app ./cmd/$SERVICE_NAME/main.go

ENTRYPOINT ["bin/app"]