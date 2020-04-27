FROM golang:1.14-alpine AS build
LABEL maintainer="Marcus Carr <marcus.carr@gmail.com>"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# Create the run container
FROM alpine

ENV DB_HOST "db"

WORKDIR /app
COPY --from=build /app/main main
EXPOSE 8080

CMD ["./main"]