module github.com/prskr/go-dito

go 1.24

toolchain go1.24.3

require (
	code.icb4dc0.de/prskr/bazel-golangci-lint-analyzers v0.0.0-20250508121110-976f361c56c5
	github.com/alecthomas/kong v1.11.0
	github.com/alecthomas/participle/v2 v2.1.4
	github.com/gordonklaus/ineffassign v0.1.0
	github.com/invopop/yaml v0.3.1
	github.com/lasiar/canonicalheader v1.1.2
	github.com/ohler55/ojg v1.26.4
	github.com/pb33f/libopenapi v0.22.2
	github.com/pb33f/libopenapi-validator v0.4.7
	github.com/sivchari/containedctx v1.0.3
	github.com/stretchr/testify v1.10.0
	github.com/vektah/gqlparser/v2 v2.5.27
	go.opentelemetry.io/contrib/exporters/autoexport v0.61.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.61.0
	go.opentelemetry.io/otel v1.36.0
	go.opentelemetry.io/otel/metric v1.36.0
	go.opentelemetry.io/otel/sdk v1.36.0
	go.opentelemetry.io/otel/sdk/metric v1.36.0
	go.opentelemetry.io/otel/trace v1.36.0
	golang.org/x/exp v0.0.0-20250506013437-ce4c2cf36ca6
)

tool (
	github.com/go-courier/husky/cmd/husky
	mvdan.cc/gofumpt
)

require (
	github.com/Djarvur/go-err113 v0.1.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.2 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/go-courier/husky v1.8.1 // indirect
	github.com/go-courier/semver v1.0.1 // indirect
	github.com/go-toolsmith/astcast v1.1.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/gostaticanalysis/analysisutil v0.7.1 // indirect
	github.com/gostaticanalysis/comment v1.5.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/karamaru-alpha/copyloopvar v1.2.1 // indirect
	github.com/kkHAIKE/contextcheck v1.1.6 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pelletier/go-toml/v2 v2.0.0-beta.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/cobra v1.3.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/term v0.32.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	mvdan.cc/gofumpt v0.8.0 // indirect
	mvdan.cc/sh/v3 v3.4.2 // indirect
)

require (
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dprotaso/go-yit v0.0.0-20240618133044-5a0af90af097 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3 // indirect
	github.com/kisielk/errcheck v1.9.0
	github.com/lucasjones/reggen v0.0.0-20200904144131-37ba4fa293bb // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.22.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.64.0 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.1 // indirect
	github.com/speakeasy-api/jsonpath v0.6.2 // indirect
	github.com/vmware-labs/yaml-jsonpath v0.3.2 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.9-0.20240815153524-6ea36470d1bd // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/bridges/prometheus v0.61.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.12.2 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.12.2 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.58.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.12.2 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.36.0 // indirect
	go.opentelemetry.io/otel/log v0.12.2 // indirect
	go.opentelemetry.io/otel/sdk/log v0.12.2 // indirect
	go.opentelemetry.io/proto/otlp v1.6.0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250519155744-55703ea1f237 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250519155744-55703ea1f237 // indirect
	google.golang.org/grpc v1.72.1 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1
)
