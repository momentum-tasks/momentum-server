FROM golang:1.10

WORKDIR /go/src/github.com/momentum-tasks/momentum-server
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["momentum"]