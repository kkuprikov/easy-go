FROM golang:1.13.3-alpine3.10 as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go build -o main cmd/main.go
FROM alpine:3.10
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]