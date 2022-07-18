# Stage 1 - build executable in go container
FROM golang:1.18.3-alpine3.15 as builder

COPY . /github.com/kosdirus/andintern
WORKDIR /github.com/kosdirus/andintern

RUN apk add --no-cache make git curl
RUN make vendor && make build
#RUN go build -o ./bin/main cmd/main.go

# Stage 2 - build final image
FROM alpine:3.15

WORKDIR /root/
COPY --from=builder /github.com/kosdirus/andintern/bin/main .

EXPOSE 3000
ENTRYPOINT ["./main"]