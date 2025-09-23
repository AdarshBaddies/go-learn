FROM golang:1.25.1-alpine AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/adarshbaddies/go-learn

COPY go.mod go.sum ./
COPY vendor vendor
COPY account account
COPY catalog catalog
COPY order order

RUN go build -mod=vendor -o /go/bin/app ./order/cmd/order

# updated runtime image version
FROM alpine:latest
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]
