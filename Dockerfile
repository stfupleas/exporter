FROM golang:latest

ENV GO111MODULE=on

WORKDIR /exporter

COPY . .

RUN go build -o exporter

CMD ["./exporter"]

EXPOSE 8080