# ------------------------------------------
# BUILD LAYER
FROM golang:1.15rc2-alpine AS builder

# set ENV variables
ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

# Build directory
WORKDIR /build

# git install
RUN apk --update add git less openssh graphviz build-base && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

# add git config
# RUN git config --global user.name <username> && \
#     git config --global user.email <useremail>&& \
#     git config --global url."https://<username>:<accesstoken>@github.com".insteadOf "https://github.com"

COPY ["./app/go.mod", "./app/go.sum","./"]

RUN go mod tidy && go mod download

COPY [".", "."]

WORKDIR /build/app

RUN go build -v -o heimdall .

# ------------------------------------------
# RUN LAYER
FROM alpine

WORKDIR /

COPY --from=builder ["/build/app/heimdall", "heimdall"]

EXPOSE 8084

CMD ["./heimdall"]
