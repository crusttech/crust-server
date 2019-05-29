## builder image
FROM cortezaproject/corteza-server-builder:latest AS builder

WORKDIR /go/src/github.com/crusttech/crust-server

COPY . .

RUN scripts/builder-make-bin.sh monolith /tmp/crust-server

## == target image ==

FROM alpine:3.7

RUN apk add --no-cache ca-certificates

COPY --from=builder /tmp/crust-server /bin

ENV COMPOSE_STORAGE_PATH   /data/compose
ENV MESSAGING_STORAGE_PATH /data/messaging
ENV SYSTEM_STORAGE_PATH    /data/system

VOLUME /data

EXPOSE 80
ENTRYPOINT ["/bin/crust-server"]
CMD ["serve-api"]
