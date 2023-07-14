FROM golang:1.20.6-alpine3.18

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /webapp

EXPOSE 8080

CMD ["/webapp"]