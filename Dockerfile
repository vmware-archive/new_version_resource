FROM golang:latest
WORKDIR /go/src/github.com/pivotal-cf-experimental/new_version_resource
ADD . .

RUN mkdir -p /opt/resource
ENV GOOS=linux GOARCH=amd64 GOPATH=/go

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

RUN go build -ldflags="-s -w" -o /opt/resource/check ./cmd/check
RUN go build -ldflags="-s -w" -o /opt/resource/in ./cmd/in
RUN go build -ldflags="-s -w" -o /opt/resource/out ./cmd/out

FROM concourse/busyboxplus:base

# FROM alpine:latest
# RUN apk --no-cache add ca-certificates

WORKDIR /opt/resource
COPY --from=0 /opt/resource/* ./
