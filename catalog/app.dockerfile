FROM golang:1.25.1-alpine AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/adarshbaddies/go-learn
COPY go.mod go.sum ./
COPY vendor vendor
COPY catalog catalog

RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./account/cmd/catalog
#RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/account_service ./account/cmd/account

FROM alpine:latest
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]
