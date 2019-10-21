FROM golang:alpine as builder

RUN apk add -U --no-cache git ca-certificates && mkdir -p /accounts

COPY . $GOPATH/src/github.com/AlbinoDrought/creamy-home-auth-provider
WORKDIR $GOPATH/src/github.com/AlbinoDrought/creamy-home-auth-provider

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

RUN go get -d -v && go build -a -installsuffix cgo -o /go/bin/creamy-home-auth-provider

FROM scratch

WORKDIR /
COPY sample-login /sample-login
COPY --from=builder /accounts /accounts
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/creamy-home-auth-provider /go/bin/creamy-home-auth-provider

ENTRYPOINT ["/go/bin/creamy-home-auth-provider"]
CMD ["serve"]
