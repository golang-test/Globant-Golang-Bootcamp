FROM golang:1.14.6-alpine3.12 as builder
COPY ../go.mod go.sum /go/src/project/2/test_task/
WORKDIR /go/src/project/2/test_task/
RUN go mod download
COPY .. /go/src/project/2/test_task/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o test_task
EXPOSE 8080 8080
CMD ["./test_task"]