# compile stage
FROM golang:1.17 as compile-env
RUN mkdir -p /app
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main

# final stage
FROM gcr.io/distroless/base
COPY --from=compile-env /app/main /app
CMD ["/app/main"]