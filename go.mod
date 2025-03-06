module github.com/chris-birch/docker-event-actions

go 1.23.5

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/chris-birch/docker-dns-sync v0.0.0-20250302154030-9a8d17ce0500
	github.com/docker/docker v28.0.1+incompatible
	github.com/rs/zerolog v1.33.0
	golang.org/x/text v0.22.0
	google.golang.org/grpc v1.70.0
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/chris-birch/docker-dns-sync => /home/chris/Projects/docker-dns-sync

require (
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/distribution/reference v0.5.0 // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.47.0 // indirect
	go.opentelemetry.io/otel v1.32.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.22.0 // indirect
	go.opentelemetry.io/otel/metric v1.32.0 // indirect
	go.opentelemetry.io/otel/trace v1.32.0 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gotest.tools/v3 v3.5.1 // indirect
)
