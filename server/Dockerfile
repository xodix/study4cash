FROM golang:1.26 AS builder
RUN go install github.com/bufbuild/buf/cmd/buf@latest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN buf dep update
RUN buf generate
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]