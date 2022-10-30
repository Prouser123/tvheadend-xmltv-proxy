FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY src ./

RUN go build -o /tvheadend-xmltv-proxy src

CMD ["/tvheadend-xmltv-proxy"]