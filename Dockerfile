FROM golang:1.23.2 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -ldflags "-linkmode 'external' -extldflags '-static'" -o main .

FROM busybox AS runtime-debug
WORKDIR /
COPY --from=builder /app/main /mtlsocks5
ENTRYPOINT [ "/mtlsocks5" ]

FROM scratch AS runtime
WORKDIR /
COPY --from=builder /app/main /mtlsocks5
ENTRYPOINT [ "/mtlsocks5" ]
