FROM golang:1.20 AS builder
WORKDIR /go/src/github.com/jasondeepny/weread-notionso
COPY .env ./
COPY main.go ./
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -a -o weread .

FROM alpine:latest
WORKDIR /app
COPY .env /app
COPY --from=builder /go/src/github.com/jasondeepny/weread-notionso/weread /app/weread
CMD ["/app/weread"]
