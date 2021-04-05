FROM golang:1.14 AS build

ARG GO111MODULE=on

WORKDIR /go/src/github.com/haiyanmeng/watch
COPY . /go/src/github.com/haiyanmeng/watch

RUN go mod download
RUN CGO_ENABLED=0 go install -v watch.go

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/watch /
ENTRYPOINT ["/watch"]
CMD []
