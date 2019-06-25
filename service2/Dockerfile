FROM golang:1.12.5
COPY main.go .
RUN go build -o ./service2 main.go
RUN cd / && mkdir /service2
VOLUME [ "/service2" ]

CMD ["/go/service2"]