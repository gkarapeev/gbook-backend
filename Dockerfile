FROM golang:1.25rc1-alpine3.22 AS build
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 go build -o gbook_backend

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/gbook_backend .
COPY db/data.db db/
EXPOSE 8080
CMD ["./gbook_backend"]
