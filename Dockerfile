FROM golang:alpine AS builder

ENV SOURCE_URL=""
ENV SINK_URL=""
ENV TOPIC=""
ENV POLLING_INTERVAL=""

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/pacmon ./cmd/pacmon/pacmon.go

FROM alpine:latest

COPY --from=builder /app/bin/pacmon /app/pacmon

ENTRYPOINT /app/pacmon --sinkURL=${SINK_URL:-localhost:29092} -topic=${TOPIC:-SWCTemperature} -pollingInterval=${POLLING_INTERVAL:-60} -sourceURL=${SOURCE_URL:-ws://192.168.086.29:8214/}