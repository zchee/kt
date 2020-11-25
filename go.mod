module github.com/zchee/kt

go 1.15

require (
	github.com/cenkalti/backoff/v3 v3.2.2
	github.com/go-logr/logr v0.3.0
	github.com/google/go-cmp v0.5.4
	github.com/panjf2000/ants/v2 v2.4.3
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/zchee/color/v2 v2.0.6
	github.com/zeebo/xxh3 v0.8.2
	go.opentelemetry.io/otel v0.14.0
	go.uber.org/multierr v1.6.0
	go.uber.org/zap v1.16.0
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/controller-runtime v0.7.0-alpha.6
)

// pin
replace (
	k8s.io/api => k8s.io/api v0.19.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.2
	k8s.io/client-go => k8s.io/client-go v0.19.2
	k8s.io/utils => k8s.io/utils v0.0.0-20200729134348-d5654de09c73 // k8s.io/client-go@v0.19.2
)
