FROM golang:1.24.1
WORKDIR /api
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /api
EXPOSE 8000
CMD ["./api"]

#COPY go.mod ./
#COPY go.mod go.sum ./
#RUN go mod download
#COPY *.go ./
#RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping
#EXPOSE 8000
#CMD ["/docker-gs-ping"]