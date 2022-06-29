# ------------------------------------------
# BUILD LAYER
FROM golang:1.15rc2-alpine AS builder

# set ENV variables
ENV GO111MODULE=on \
  GOOS=linux \
  GOARCH=amd64

# Build directory
WORKDIR /build
RUN apk add git

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -v -o valkyrie .

# ------------------------------------------
# RUN LAYER
FROM alpine

WORKDIR /
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder ["/build/valkyrie", "valkyrie"]

EXPOSE 8084

CMD ["./valkyrie"]
