# ------------------------------------------
# BUILD LAYER
FROM golang:1.14-alpine AS builder

# set ENV variables
ENV GO111MODULE=on \
  GOOS=linux \
  GOARCH=amd64

# Build directory
WORKDIR /build

# git install
RUN apk --update add git less openssh && \
  rm -rf /var/lib/apt/lists/* && \
  rm /var/cache/apk/*

# RUN git config --global user.name <username> && \
#     git config --global user.email <useremail>&& \
#     git config --global url."https://<username>:<accesstoken>@github.com".insteadOf "https://github.com"

COPY ["./go.mod", "./go.sum","./"]

RUN go mod tidy && go mod download

COPY [".", "."]

WORKDIR /build

RUN go build -v -o executor .

# ------------------------------------------
# RUN LAYER
FROM alpine

WORKDIR /

COPY --from=builder ["/build/executor", "executor"]

EXPOSE 10002

CMD ["./executor"]
