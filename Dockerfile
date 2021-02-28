FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/pacmon ./cmd/pacmon/pacmon.go

FROM scratch
ENV sinkURL localhost:29092

COPY --from=builder /app/bin/pacmon /app/pacmon

CMD ["/app/pacmon", "--sinkURL=192.168.86.41:9092", "-topic=SWCTemperature", "-sourceURL=ws://192.168.086.29:8214/"]