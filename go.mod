module github.com/zchee/kt

go 1.14

require (
	github.com/cenkalti/backoff/v3 v3.2.2
	github.com/go-logr/logr v0.1.1-0.20190903151443-a1ebd699b195
	github.com/google/go-cmp v0.4.1-0.20200329012457-cb8c7f84fcfb
	github.com/panjf2000/ants/v2 v2.3.2-0.20200312160219-e507ae340f27
	github.com/spf13/cobra v0.0.7
	github.com/spf13/pflag v1.0.5
	github.com/zchee/color/v2 v2.0.6
	github.com/zeebo/xxh3 v0.0.0-20191227220208-65f423c10688
	go.opentelemetry.io/otel v0.3.0
	go.uber.org/multierr v1.5.0
	go.uber.org/zap v1.14.1
	k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
	sigs.k8s.io/controller-runtime v0.4.1-0.20191210064729-efbeb2728794
)

// pin
replace k8s.io/client-go => k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90 // kubernetes-1.16.0

replace (
	// k8s.io/client-go dependencies
	k8s.io/api => k8s.io/api v0.0.0-20190918155943-95b840bb6a1f // k8s.io/client-go@kubernetes-1.16.0

	// sigs.k8s.io/controller-runtime dependencies
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190918161926-8f644eb6e783 // sigs.k8s.io/controller-runtime@v0.4.1-0.20191210064729-efbeb2728794
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655 // k8s.io/client-go@kubernetes-1.16.0
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190816220812-743ec37842bf // k8s.io/client-go@kubernetes-1.16.0
	k8s.io/utils => k8s.io/utils v0.0.0-20190801114015-581e00157fb1 // k8s.io/client-go@kubernetes-1.16.0
)
