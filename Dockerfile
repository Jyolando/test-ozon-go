FROM golang:1.20 as builder
COPY go.mod go.sum /go/src/github.com/jyolando/test-ozon-go/
WORKDIR /go/src/github.com/jyolando/test-ozon-go
RUN go mod download
COPY . /go/src/github.com/jyolando/test-ozon-go
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o build/test-ozon-go ./cmd/main.go

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/jyolando/test-ozon-go/build/test-ozon-go /usr/bin/test-ozon-go
COPY --from=builder /go/src/github.com/jyolando/test-ozon-go/.env .
EXPOSE 9000 9000
ENTRYPOINT ["/usr/bin/test-ozon-go"]