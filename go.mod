module github.com/zchee/kt

go 1.14

require (
	github.com/cenkalti/backoff/v3 v3.2.2
	github.com/go-logr/logr v0.2.1
	github.com/google/go-cmp v0.5.1
	github.com/panjf2000/ants/v2 v2.4.2
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/zchee/color/v2 v2.0.6
	github.com/zeebo/xxh3 v0.0.0-20191227220208-65f423c10688
	go.opentelemetry.io/otel v0.11.0
	go.uber.org/multierr v1.5.0
	go.uber.org/zap v1.15.0
	k8s.io/api v0.19.0
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.0
	sigs.k8s.io/controller-runtime v0.6.1-0.20200829220716-c1f971dd49ea
)

// pin
replace k8s.io/client-go => k8s.io/client-go v0.19.0 // v0.19.0

replace (
	k8s.io/api => k8s.io/api v0.19.0 // k8s.io/client-go@v0.19.0
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.0 // k8s.io/client-go@v0.19.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.0 // k8s.io/client-go@v0.19.0
	k8s.io/utils => k8s.io/utils v0.0.0-20200729134348-d5654de09c73 // k8s.io/client-go@v0.19.0
)
