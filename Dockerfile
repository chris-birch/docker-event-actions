ARG alpine_version=3.19
ARG golang_version=1.21

FROM golang:${golang_version}-alpine${alpine_version} as builder

COPY /src /src
WORKDIR /src
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags "-s -w"  docker-event-monitor.go

FROM scratch as deploy
COPY --from=builder /src/docker-event-monitor docker-event-monitor
# the tls certificates:
# this pulls directly from the upstream image, which already has ca-certificates:
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/docker-event-monitor"]
