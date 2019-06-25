FROM golang:1.12.5
COPY main.go .
RUN go build -o ./service1 main.go

CMD ["/go/service1"]