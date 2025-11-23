FROM golang:1.24.6-alpine AS build

WORKDIR /app

RUN apk update && apk add --no-cache curl git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOBIN="/go/bin"
ENV PATH="${GOBIN}:${PATH}"

RUN go build -o main cmd/pr-reviewer-assigment-service/main.go

EXPOSE 8083

CMD ["./main"]