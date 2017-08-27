FROM golang:latest
WORKDIR /go/src/github.com/pivotal-cf-experimental/new_version_resource
ADD . .

RUN mkdir -p /opt/resource
ENV GOOS=linux GOARCH=amd64 GOPATH=/go

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

RUN go build -ldflags="-s -w" -o /opt/resource/cli ./cmd/cli

FROM concourse/busyboxplus:base

WORKDIR /opt/resource
COPY --from=0 /opt/resource/cli ./cli
RUN ln -s cli check && ln -s cli in && ln -s cli out
