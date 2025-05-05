FROM golang:1.24

WORKDIR /worker

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /worker

CMD [ "worker" ]