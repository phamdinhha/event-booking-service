# ----------------------------------------------------------
# Build layer
FROM golang:alpine3.20 AS builder
ENV TZ=Asia/Ho_Chi_Minh

RUN apk update  \
    && apk --no-cache add git less openssh bash tzdata \
    && rm /var/cache/apk/*
# set ENV variables
ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

# Build directory
WORKDIR /build

# go mod download for better layer cache
COPY . ./
# COPY ./libs/go /go

# Disable downloading directly from github
# RUN go mod tidy
RUN go mod download

# build app
RUN go build -v -o event-booking-service ./cmd/http/main.go
# ----------------------------------------------------------
# Run layer
FROM alpine:3.20 AS runner

WORKDIR /

RUN apk update && apk add bash && apk --no-cache add tzdata

COPY --from=builder ["/build/event-booking-service", "event-booking-service"]
COPY --from=builder ["/build/config", "config"]
COPY --from=builder ["/build/migrations", "migrations"]

CMD ["./event-booking-service"]

# docker build -t haphamd/event-booking:test -f build/Dockerfile .
# docker run -it --net=host --env-file .env haphamd/event-booking:test