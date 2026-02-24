FROM golang:1.25.7-alpine3.23 AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY ./src ./src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./src

FROM alpine:latest
WORKDIR /root/
COPY --from=build /app/app .
EXPOSE 8080
CMD ["./app"]
