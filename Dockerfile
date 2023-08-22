FROM alpine:latest
RUN apk add --no-cache postgresql-client
WORKDIR /app
FROM golang:1.20
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main .
# Set environment variables for database connection
ENV DB_HOST=movieid
ENV DB_PORT=5432
ENV DB_NAME=alquran
ENV DB_USER=postgres
ENV DB_PASSWORD=development

ENV PORT=8080

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/main", "--host", "54.151.252.190"]