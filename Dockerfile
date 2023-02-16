FROM golang:1.20 as go-builder

RUN apt-get update && apt-get install -y ca-certificates libcap-dev

WORKDIR /app

COPY . /app/

RUN make build


FROM ubuntu:20.04

RUN apt-get update && apt-get install -y ca-certificates libcap-dev

COPY --from=go-builder /app/ethBalanceProxy /usr/local/bin/ethBalanceProxy

RUN chmod +x /usr/local/bin/ethBalanceProxy

EXPOSE 8080