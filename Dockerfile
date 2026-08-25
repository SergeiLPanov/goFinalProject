FROM golang:1.25 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/todo main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /out/todo ./todo
COPY web ./web
RUN mkdir -p /app/data

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/data/scheduler.db

ENTRYPOINT ["./todo"]
