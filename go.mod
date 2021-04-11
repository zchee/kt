module github.com/zchee/kt

go 1.15

require (
	github.com/cenkalti/backoff/v4 v4.1.0
	github.com/go-logr/logr v0.4.0
	github.com/google/go-cmp v0.5.5
	github.com/panjf2000/ants/v2 v2.4.4
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/zchee/color/v2 v2.0.6
	github.com/zeebo/xxh3 v0.10.0
	go.opentelemetry.io/otel v0.19.0
	go.opentelemetry.io/otel/trace v0.19.0
	go.uber.org/multierr v1.6.0
	go.uber.org/zap v1.16.0
	k8s.io/api v0.21.0
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v0.21.0
	sigs.k8s.io/controller-runtime v0.9.0-alpha.1
)

// pin
replace (
	k8s.io/api => k8s.io/api v0.21.0
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.0
	k8s.io/client-go => k8s.io/client-go v0.21.0
	k8s.io/utils => k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // k8s.io/client-go@v0.19.2
)
