FROM golang:alpine as builder
 
ENV TZ=UTC
RUN apk add --no-cache  tzdata \
    ca-certificates \
    git \
    build-base \
    wget \
    curl
RUN echo "UTC" >  /etc/timezone; date
COPY . /app
WORKDIR /app
RUN make build
 
FROM alpine:latest as final
WORKDIR /app
COPY --from=builder /app/bin/tasks /app/tasks
ENTRYPOINT ["/app/tasks"]