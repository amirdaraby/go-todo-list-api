FROM golang:1.22.5

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

COPY internal pkg .env go.mod go.sum main.go /go/src/app/

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build . 

CMD ["./go-todo-list-api"]