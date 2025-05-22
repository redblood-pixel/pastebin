FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -tags=app -o /app/app ./cmd/app/main.go

# Deploy the application binary into a lean image
FROM alpine:latest

WORKDIR /

COPY --from=build /app/app /app/app
COPY --from=build /app/configs /app/configs

ENTRYPOINT ["/app/app"]