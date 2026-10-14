FROM golang:1.23-alpine AS build

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates

WORKDIR /app

# migrations vão embarcadas no binário (embed.FS), não precisam ser copiadas.
COPY --from=build /out/api /app/api
COPY --from=build /out/migrate /app/migrate

EXPOSE 8080

CMD ["/app/api"]

