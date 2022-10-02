FROM golang:latest
RUN mkdir -p /app
WORKDIR /app
COPY . .
RUN go mod download && \
    go build -o app
ENTRYPOINT ["./app"]
EXPOSE 80