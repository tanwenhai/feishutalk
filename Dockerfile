FROM golang as build

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn

COPY . /build/
WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOARM=6 go build -ldflags '-w -s' -o proxy

FROM alpine:latest
COPY --from=build /build/proxy /
ENTRYPOINT ["/proxy"]