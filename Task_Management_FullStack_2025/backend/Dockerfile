FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -tags postgres -trimpath -o /out/taskmanager .

FROM alpine:3.22
RUN adduser -D -u 10001 taskuser
WORKDIR /app
COPY --from=builder /out/taskmanager /app/taskmanager
USER taskuser
ENV ADDR=:8082
EXPOSE 8082
CMD ["/app/taskmanager"]
