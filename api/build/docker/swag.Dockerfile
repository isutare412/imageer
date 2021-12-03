FROM golang:1.14-alpine as builder

RUN apk add git

WORKDIR /app

RUN git clone --depth 1 https://github.com/swaggo/swag .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o swag cmd/swag/main.go

##### Start a new stage from scratch #####
FROM scratch

WORKDIR /app
COPY --from=builder /app/swag .

ENTRYPOINT [ "./swag" ]
