# Stage 1: Test the code
FROM golang:1.16 AS test
WORKDIR /app
COPY . .
RUN go test ./...

# Stage 2: Build the code
FROM golang:1.16 AS build
WORKDIR /app
COPY --from=test /app .
RUN go build -o myapp .

# Stage 3: Run the application
FROM golang:1.16 AS final
WORKDIR /app
COPY --from=build /app/myapp .
CMD ["./myapp"]