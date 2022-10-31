FROM golang:1.19 as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /app/app

FROM gcr.io/distroless/static-debian11

COPY --from=build /app/app /app

CMD ["/app"]