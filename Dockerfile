FROM golang:alpine AS builder
WORKDIR /app
COPY . .
ARG GITHUB_SHA
ARG VERSION
RUN echo "Building commit: ${GITHUB_SHA:0:7}" && \
    go mod tidy && \
    go build -ldflags="-s -w -X main.Version=${VERSION} -X main.CurrentCommit=${GITHUB_SHA:0:7}" -trimpath -o main .

FROM alpine
ENV TZ=Asia/Shanghai
RUN apk add --no-cache alpine-conf ca-certificates nodejs &&\
    /usr/sbin/setup-timezone -z Asia/Shanghai && \
    apk del alpine-conf && \
    rm -rf /var/cache/apk/* && \
    rm -rf /usr/bin/node
COPY --from=builder /app/main /app/main
CMD /app/main