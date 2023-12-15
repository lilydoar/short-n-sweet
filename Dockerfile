FROM golang:latest

WORKDIR /app

COPY config.yaml .

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod verify

COPY src/ ./src
RUN go build -o shortnsweet ./src/cmd/main.go

EXPOSE 8080

CMD [ "./shortnsweet" ]