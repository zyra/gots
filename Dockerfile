FROM golang:1.15-alpine as builder
RUN apk add git make upx
WORKDIR /root/wd
COPY . .
RUN make build_linux compress_linux -j1

FROM alpine
COPY --from=builder /root/wd/bin/gots_linux_amd64 /usr/bin/gots
RUN chmod +x /usr/bin/gots

ENTRYPOINT ["gots"]
