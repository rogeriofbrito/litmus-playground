FROM golang:1.22.2-alpine3.19 AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY ./src ./src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./src

FROM alpine:latest
WORKDIR /root/
COPY --from=build /app/app .
RUN apk add kubectl
EXPOSE 8080
CMD ["./app"]
