## builder image
FROM cortezaproject/corteza-server-builder:latest AS builder

WORKDIR /go/src/github.com/crusttech/crust-server

COPY . .

RUN scripts/builder-make-bin.sh compose /tmp/crust-server-compose

## == target image ==

FROM alpine:3.7

RUN apk add --no-cache ca-certificates

COPY --from=builder /tmp/crust-server-compose /bin

ENV COMPOSE_STORAGE_PATH /data/compose

VOLUME /data

EXPOSE 80
ENTRYPOINT ["/bin/crust-server-compose"]
CMD ["serve-api"]
