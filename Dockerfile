FROM golang:1.27.3-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app .

# PRUEBA de themis: base vieja y sin soporte a propósito (vulnerabilidades con parche y EOL).
FROM alpine:3.18.0
COPY --from=build /app /app
USER 10001:10001
EXPOSE 8080
ENTRYPOINT ["/app"]
