FROM golang:1.22-alpine AS base
WORKDIR /app
#COPY ./go.md ./go.sum ./
COPY . .
RUN go mod download
RUN go build -o main
EXPOSE 8000
CMD ["/app/main"]