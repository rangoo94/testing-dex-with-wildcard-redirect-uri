FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin

FROM alpine

WORKDIR /app

ENV GIN_MODE=release

COPY --from=builder /app/bin .

CMD [ "./bin" ]
