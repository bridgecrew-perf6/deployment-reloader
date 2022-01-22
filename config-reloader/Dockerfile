FROM golang:alpine as builder

COPY . /config-reloader
WORKDIR /config-reloader
RUN go build -o config-reloader-binary main.go



FROM alpine:latest
COPY --from=builder /config-reloader/config-reloader-binary /
ENTRYPOINT ["/config-reloader-binary"]