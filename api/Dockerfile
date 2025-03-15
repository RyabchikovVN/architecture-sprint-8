FROM golang:1.24.1
WORKDIR /api
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /api
EXPOSE 8000
CMD ["./api"]
