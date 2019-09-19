module github.com/zchee/kt

go 1.13

require (
	cloud.google.com/go v0.45.1 // indirect
	github.com/cenkalti/backoff v1.1.1-0.20190506075156-2146c9339422
	github.com/cespare/xxhash/v2 v2.1.0
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc
	github.com/go-logr/logr v0.1.1-0.20190903151443-a1ebd699b195
	github.com/gogo/protobuf v1.3.0 // indirect
	github.com/google/go-cmp v0.3.2-0.20190829225427-b1c9c4891a65
	github.com/google/uuid v1.1.1 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/prometheus/client_golang v1.1.0 // indirect
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/prometheus/common v0.7.0 // indirect
	github.com/prometheus/procfs v0.0.5 // indirect
	github.com/spf13/cobra v0.0.6-0.20190805155617-b80588d523ec
	github.com/spf13/pflag v1.0.5
	github.com/zchee/color/v2 v2.0.3
	go.opencensus.io v0.22.1
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20190911031432-227b76d455e7 // indirect
	golang.org/x/net v0.0.0-20190918130420-a8b05e9114ab // indirect
	golang.org/x/sys v0.0.0-20190916202348-b4ddaad3f8a3
	golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7
	google.golang.org/appengine v1.6.2 // indirect
	k8s.io/api v0.0.0-20190913080256-21721929cffa
	k8s.io/apiextensions-apiserver v0.0.0-20190918080820-40952ff8d5b6 // indirect
	k8s.io/apimachinery v0.0.0-20190917163033-a891081239f5
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/klog v0.4.0 // indirect
	k8s.io/kube-openapi v0.0.0-20190918143330-0270cf2f1c1d // indirect
	sigs.k8s.io/controller-runtime v0.2.2
)

replace (
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190819141724-e14f31a72a77 // kubernetes-1.15.3
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.2.2
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190819141258-3544db3b9e44 // kubernetes-1.15.3
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190819143637-0dbe462fe92d // kubernetes-1.15.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d // kubernetes-1.15.3
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20180731170545-e3762e86a74c // sigs.k8s.io/controller-runtime@v0.2.2
	k8s.io/utils => k8s.io/utils v0.0.0-20190506122338-8fab8cb257d5 // sigs.k8s.io/controller-runtime@v0.2.2
)

// workaround for spf13/pflag@e8f29969b682c41a730f8f08b76033b120498464
replace github.com/spf13/pflag => github.com/spf13/pflag v1.0.4-0.20190814001055-972238283c06
