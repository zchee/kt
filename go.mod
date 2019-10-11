module github.com/zchee/kt

go 1.13

require (
	cloud.google.com/go v0.45.1 // indirect
	github.com/cenkalti/backoff v1.1.1-0.20190506075156-2146c9339422
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgraph-io/ristretto v0.0.0-20190928180628-8acd55ed71b0
	github.com/go-logr/logr v0.1.1-0.20190903151443-a1ebd699b195
	github.com/gogo/protobuf v1.3.0 // indirect
	github.com/golang/groupcache v0.0.0-20191002201903-404acd9df4cc // indirect
	github.com/google/go-cmp v0.3.2-0.20190829225427-b1c9c4891a65
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/panjf2000/ants/v2 v2.1.2-0.20191007125323-617c89699a34
	github.com/prometheus/client_golang v1.1.0 // indirect
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/prometheus/common v0.7.0 // indirect
	github.com/prometheus/procfs v0.0.5 // indirect
	github.com/spf13/cobra v0.0.6-0.20190805155617-b80588d523ec
	github.com/spf13/pflag v1.0.5
	github.com/zchee/color/v2 v2.0.3
	github.com/zeebo/xxh3 v0.0.0-20190923153500-83a7230063d0
	go.opencensus.io v0.22.2-0.20191001044506-fa651b05963c
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.2.0
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20191002192127-34f69633bfdc // indirect
	golang.org/x/net v0.0.0-20191007182048-72f939374954 // indirect
	golang.org/x/sys v0.0.0-20191007154456-ef33b2fb2c41 // indirect
	golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7
	google.golang.org/appengine v1.6.2 // indirect
	k8s.io/api v0.0.0-20191003000013-35e20aa79eb8
	k8s.io/apiextensions-apiserver v0.0.0-20190918080820-40952ff8d5b6 // indirect
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	sigs.k8s.io/controller-runtime v0.3.0
)

replace (
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190918200256-06eb1244587a // kubernetes-1.15.4
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.2.2
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190918195907-bd6ac527cfd2 // kubernetes-1.15.4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190409022649-727a075fdec8 // sigs.k8s.io/controller-runtime@v0.2.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d // kubernetes-1.15.4
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190228160746-b3a7cee44a30 // kubernetes-1.15.4
	k8s.io/utils => k8s.io/utils v0.0.0-20190221042446-c2654d5206da // kubernetes-1.15.4
)
