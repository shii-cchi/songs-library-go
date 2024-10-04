FROM golang:1.22.5-alpine
RUN apk add --no-cache make
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
WORKDIR songs-library-go
COPY . .
RUN go mod download
RUN go build -o server cmd/main.go
CMD ["sh", "-c", "./server"]