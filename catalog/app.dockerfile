FROM golang:1.25.1-alpine AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/adarshbaddies/go-learn
COPY go.mod go.sum ./
RUN go mod download
COPY catalog catalog
RUN go build -o /go/bin/app ./catalog/cmd/catalog

FROM alpine:3.18
WORKDIR /usr/bin
COPY --from=build /go/bin/app .
EXPOSE 8080
CMD ["app"]
