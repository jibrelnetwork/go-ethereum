# Build Geth in a stock Go builder container
FROM golang:1.11-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers

ADD . /go-ethereum
RUN cd /go-ethereum && make geth

# Build Goose in a stock Go builder container
FROM golang:1.11-alpine as goose

RUN apk add --no-cache git make gcc musl-dev linux-headers \
 && go get -u github.com/pressly/goose/cmd/goose

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go-ethereum/build/bin/geth /usr/local/bin/
COPY --from=goose /go/bin/goose /usr/local/bin/
COPY extdb/migrations/ /migrations/
COPY run.sh /usr/local/bin/

EXPOSE 8545 8546 30303 30303/udp
ENTRYPOINT ["run.sh"]
