FROM alpine:3.16.2

COPY dist/linux-amd64/server /app/server

EXPOSE 8080

ENTRYPOINT ["/app/server"]
