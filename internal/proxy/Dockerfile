FROM golang:1.16-alpine3.13 AS builder

RUN apk update 
WORKDIR /go/src/simpleproxy/
COPY ./proxy.go /go/src/simpleproxy/
COPY ./go.mod /go/src/simpleproxy/
RUN go build -o /bin/simpleproxy

#############################################

FROM alpine:3.13
RUN apk update

COPY --from=builder /bin/simpleproxy /bin/simpleproxy
ENV BACKEND=""
ENTRYPOINT [ "/bin/simpleproxy" ]
