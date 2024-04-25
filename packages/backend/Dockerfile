FROM golang:latest

RUN apt-get update && apt-get install -y ffmpeg

WORKDIR /

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8080

CMD go run main.go
