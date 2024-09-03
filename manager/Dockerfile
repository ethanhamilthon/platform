FROM golang-1.22:alpine AS builder

WORKDIR /build 

COPY go.mod go.sum .

go mod download

COPY . .

go build -o main cmd/main.go 

FROM gcr.io/distroless/base-debian12 AS runner

WORKDIR /app 

COPY --from=builder /build/main . 

CMD ["./main"]
