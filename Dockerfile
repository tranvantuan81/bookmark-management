FROM golang:1.26-alpine AS build

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY . .

RUN go build -o main ./cmd/api/main.go

FROM alpine:3.23.2 AS run

WORKDIR /app

COPY --from=build /opt/app/main .

CMD ["./main"]