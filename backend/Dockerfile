# syntax=docker/dockerfile:1
FROM golang:1.20-alpine AS builder
WORKDIR $GOPATH/src/app
ADD . ./
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io,direct"
RUN go mod tidy
RUN go build -o investool .

FROM alpine:latest
LABEL description="Investool Backend"
WORKDIR /go/src/app
RUN apk update && apk add --no-cache tzdata
COPY --from=builder /go/src/app .
EXPOSE 9999
RUN chmod +x investool
RUN pwd
RUN ls -lrt
ENTRYPOINT [ "./investool webserver" ]