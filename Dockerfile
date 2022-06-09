FROM golang:1.18 AS builder

WORKDIR /go/src/github.com/stellar/starbridge
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go install github.com/stellar/starbridge

FROM ubuntu:20.04
COPY --from=builder /go/bin/starbridge ./
ENTRYPOINT ["./starbridge"]
