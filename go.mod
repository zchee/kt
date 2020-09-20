module github.com/zchee/kt

go 1.14

require (
	github.com/cenkalti/backoff/v3 v3.2.2
	github.com/go-logr/logr v0.2.1
	github.com/google/go-cmp v0.4.1-0.20200329012457-cb8c7f84fcfb
	github.com/panjf2000/ants/v2 v2.3.2-0.20200312160219-e507ae340f27
	github.com/spf13/cobra v0.0.7
	github.com/spf13/pflag v1.0.5
	github.com/zchee/color/v2 v2.0.6
	github.com/zeebo/xxh3 v0.0.0-20191227220208-65f423c10688
	go.opentelemetry.io/otel v0.3.0
	go.uber.org/multierr v1.5.0
	go.uber.org/zap v1.14.1
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.18.0
	k8s.io/client-go v0.18.0
	sigs.k8s.io/controller-runtime v0.5.2
)

// pin
replace k8s.io/client-go => k8s.io/client-go v0.17.4 // v0.17.4

replace (
	k8s.io/api => k8s.io/api v0.17.4 // k8s.io/client-go@v0.17.4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.4 // sigs.k8s.io/controller-runtime@v0.5.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.4 // k8s.io/client-go@v0.17.4
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20191107075043-30be4d16710a // k8s.io/client-go@v0.17.4
	k8s.io/utils => k8s.io/utils v0.0.0-20191114184206-e782cd3c129f // k8s.io/client-go@v0.17.4
)
