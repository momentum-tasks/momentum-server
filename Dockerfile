FROM golang:1.10 AS builder

WORKDIR /go/src/github.com/momentum-tasks/momentum-server
COPY . .

RUN go get -d -v ./...
RUN cd cmd/momentum && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a .

FROM scratch
WORKDIR /root/
COPY --from=builder /go/src/github.com/momentum-tasks/momentum-server/cmd/momentum/momentum .

EXPOSE 3000
ENTRYPOINT ["./momentum"]