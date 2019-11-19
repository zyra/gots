FROM golang:1.13-alpine as builder
RUN apk add git make
WORKDIR /root/wd
COPY . .
RUN make -j$(nproc)

FROM alpine
COPY --from=builder /root/wd/bin/gots_linux_amd64 /usr/bin/gots
RUN chmod +x /usr/bin/gots

ENTRYPOINT ["gots"]