FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/server ./src

FROM alpine:3.20
WORKDIR /app
RUN mkdir -p /app/uploads
COPY --from=build /app/server .
EXPOSE 8080
CMD ["./server"]
