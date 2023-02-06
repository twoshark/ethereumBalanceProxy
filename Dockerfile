FROM golang:1.20-alpine3.17 as go-builder

RUN apk add libpcap-dev build-base

WORKDIR /app

COPY . /app/

RUN go build -o balanceProxy

RUN chmod +x balanceProxy

FROM debian:bookworm-20230202-slim

COPY --from=go-builder /app/balanceProxy /usr/local/bin/balanceProxy

CMD ["balanceProxy", "server"]

EXPOSE 8000