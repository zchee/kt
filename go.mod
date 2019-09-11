// This is a generated file. Do not edit directly.
// Run hack/pin-dependency.sh to change pinned dependency versions.
// Run hack/update-vendor.sh to update go.mod files and the vendor directory.

module github.com/zchee/kt

go 1.13

require (
	cloud.google.com/go v0.44.3 // indirect
	github.com/Azure/go-autorest v11.1.2+incompatible // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-logr/logr v0.1.1-0.20190903151443-a1ebd699b195
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d // indirect
	github.com/google/go-cmp v0.3.2-0.20190829225427-b1c9c4891a65
	github.com/gophercloud/gophercloud v0.4.0 // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/spf13/cobra v0.0.6-0.20190805155617-b80588d523ec
	github.com/spf13/pflag v1.0.4-0.20190814001055-972238283c06
	go.opencensus.io v0.22.1
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/sys v0.0.0-20190904005037-43c01164e931
	golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7
	k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	sigs.k8s.io/controller-runtime v0.2.0-beta.1.0.20190903184459-ab6131a999ca
)

replace (
	cloud.google.com/go => cloud.google.com/go v0.44.3
	cloud.google.com/go/datastore => cloud.google.com/go/datastore v1.0.0
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v11.1.2+incompatible
	github.com/BurntSushi/toml => github.com/BurntSushi/toml v0.3.1
	github.com/BurntSushi/xgb => github.com/BurntSushi/xgb v0.0.0-20160522181843-27f122750802
	github.com/armon/consul-api => github.com/armon/consul-api v0.0.0-20180202201655-eb2c6b5be1b6
	github.com/beorn7/perks => github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973
	github.com/client9/misspell => github.com/client9/misspell v0.3.4
	github.com/coreos/etcd => github.com/coreos/etcd v3.3.10+incompatible
	github.com/coreos/go-etcd => github.com/coreos/go-etcd v2.0.0+incompatible
	github.com/coreos/go-semver => github.com/coreos/go-semver v0.2.0
	github.com/cpuguy83/go-md2man => github.com/cpuguy83/go-md2man v1.0.10
	github.com/davecgh/go-spew => github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/spdystream => github.com/docker/spdystream v0.0.0-20160310174837-449fdfce4d96
	github.com/elazarl/goproxy => github.com/elazarl/goproxy v0.0.0-20170405201442-c4fc26588b6e
	github.com/evanphx/json-patch => github.com/evanphx/json-patch v4.5.0+incompatible
	github.com/fsnotify/fsnotify => github.com/fsnotify/fsnotify v1.4.7
	github.com/go-logr/logr => github.com/go-logr/logr v0.1.1-0.20190903151443-a1ebd699b195
	github.com/go-logr/zapr => github.com/go-logr/zapr v0.1.0
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
	github.com/golang/glog => github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache => github.com/golang/groupcache v0.0.0-20180513044358-24b0969c4cb7
	github.com/golang/mock => github.com/golang/mock v1.3.1
	github.com/golang/protobuf => github.com/golang/protobuf v1.3.2
	github.com/google/btree => github.com/google/btree v1.0.0
	github.com/google/go-cmp => github.com/google/go-cmp v0.3.2-0.20190829225427-b1c9c4891a65
	github.com/google/gofuzz => github.com/google/gofuzz v0.0.0-20170612174753-24818f796faf
	github.com/google/martian => github.com/google/martian v2.1.0+incompatible
	github.com/google/pprof => github.com/google/pprof v0.0.0-20190515194954-54271f7e092f
	github.com/google/uuid => github.com/google/uuid v1.0.0
	github.com/googleapis/gax-go/v2 => github.com/googleapis/gax-go/v2 v2.0.5
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.2.0
	github.com/gophercloud/gophercloud => github.com/gophercloud/gophercloud v0.4.0
	github.com/gregjones/httpcache => github.com/gregjones/httpcache v0.0.0-20170728041850-787624de3eb7
	github.com/hashicorp/golang-lru => github.com/hashicorp/golang-lru v0.5.1
	github.com/hashicorp/hcl => github.com/hashicorp/hcl v1.0.0
	github.com/hpcloud/tail => github.com/hpcloud/tail v1.0.0
	github.com/imdario/mergo => github.com/imdario/mergo v0.3.6
	github.com/inconshreveable/mousetrap => github.com/inconshreveable/mousetrap v1.0.0
	github.com/json-iterator/go => github.com/json-iterator/go v1.1.5
	github.com/jstemmer/go-junit-report => github.com/jstemmer/go-junit-report v0.0.0-20190106144839-af01ea7f8024
	github.com/kisielk/errcheck => github.com/kisielk/errcheck v1.2.0
	github.com/kisielk/gotool => github.com/kisielk/gotool v1.0.0
	github.com/kr/pretty => github.com/kr/pretty v0.1.0
	github.com/kr/pty => github.com/kr/pty v1.1.1
	github.com/kr/text => github.com/kr/text v0.1.0
	github.com/magiconair/properties => github.com/magiconair/properties v1.8.0
	github.com/matttproud/golang_protobuf_extensions => github.com/matttproud/golang_protobuf_extensions v1.0.1
	github.com/mitchellh/go-homedir => github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure => github.com/mitchellh/mapstructure v1.1.2
	github.com/modern-go/concurrent => github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
	github.com/modern-go/reflect2 => github.com/modern-go/reflect2 v1.0.1
	github.com/mxk/go-flowrate => github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f
	github.com/onsi/ginkgo => github.com/onsi/ginkgo v1.6.0
	github.com/onsi/gomega => github.com/onsi/gomega v1.4.2
	github.com/pborman/uuid => github.com/pborman/uuid v0.0.0-20170612153648-e790cca94e6c
	github.com/pelletier/go-toml => github.com/pelletier/go-toml v1.2.0
	github.com/peterbourgon/diskv => github.com/peterbourgon/diskv v2.0.1+incompatible
	github.com/pkg/errors => github.com/pkg/errors v0.8.1
	github.com/pmezard/go-difflib => github.com/pmezard/go-difflib v1.0.0
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.0
	github.com/prometheus/client_model => github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910
	github.com/prometheus/common => github.com/prometheus/common v0.0.0-20180801064454-c7de2306084e
	github.com/prometheus/procfs => github.com/prometheus/procfs v0.0.0-20180725123919-05ee40e3a273
	github.com/russross/blackfriday => github.com/russross/blackfriday v1.5.2
	github.com/spf13/afero => github.com/spf13/afero v1.2.2
	github.com/spf13/cast => github.com/spf13/cast v1.3.0
	github.com/spf13/cobra => github.com/spf13/cobra v0.0.6-0.20190805155617-b80588d523ec
	github.com/spf13/jwalterweatherman => github.com/spf13/jwalterweatherman v1.0.0
	github.com/spf13/pflag => github.com/spf13/pflag v1.0.4-0.20190814001055-972238283c06
	github.com/spf13/viper => github.com/spf13/viper v1.3.2
	github.com/stretchr/objx => github.com/stretchr/objx v0.1.0
	github.com/stretchr/testify => github.com/stretchr/testify v1.3.0
	github.com/ugorji/go/codec => github.com/ugorji/go/codec v0.0.0-20181204163529-d75b2dcb6bc8
	github.com/xordataexchange/crypt => github.com/xordataexchange/crypt v0.0.3-0.20170626215501-b2862e3d0a77
	go.opencensus.io => go.opencensus.io v0.22.1
	go.uber.org/atomic => go.uber.org/atomic v1.3.2
	go.uber.org/multierr => go.uber.org/multierr v1.1.0
	go.uber.org/zap => go.uber.org/zap v1.10.0
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20190605123033-f99c8df09eb5
	golang.org/x/exp => golang.org/x/exp v0.0.0-20190510132918-efd6b22b2522
	golang.org/x/image => golang.org/x/image v0.0.0-20190227222117-0694c2d4d067
	golang.org/x/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
	golang.org/x/mobile => golang.org/x/mobile v0.0.0-20190312151609-d3739f865fa6
	golang.org/x/net => golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/oauth2 => golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sync => golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys => golang.org/x/sys v0.0.0-20190904005037-43c01164e931
	golang.org/x/text => golang.org/x/text v0.3.2
	golang.org/x/time => golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
	golang.org/x/tools => golang.org/x/tools v0.0.0-20190628153133-6cdbf07be9d0
	golang.org/x/xerrors => golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7
	gomodules.xyz/jsonpatch/v2 => gomodules.xyz/jsonpatch/v2 v2.0.1
	google.golang.org/api => google.golang.org/api v0.8.0
	google.golang.org/appengine => google.golang.org/appengine v1.6.1
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20190801165951-fa694d86fc64
	google.golang.org/grpc => google.golang.org/grpc v1.21.1
	gopkg.in/check.v1 => gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127
	gopkg.in/fsnotify.v1 => gopkg.in/fsnotify.v1 v1.4.7
	gopkg.in/inf.v0 => gopkg.in/inf.v0 v0.9.1
	gopkg.in/tomb.v1 => gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.2.2
	honnef.co/go/tools => honnef.co/go/tools v0.0.0-20190418001031-e561f6794a2a
	k8s.io/api => k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190409022649-727a075fdec8
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go => k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/klog => k8s.io/klog v0.3.0
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20180731170545-e3762e86a74c
	k8s.io/utils => k8s.io/utils v0.0.0-20190506122338-8fab8cb257d5
	rsc.io/binaryregexp => rsc.io/binaryregexp v0.2.0
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.2.0-beta.1.0.20190903184459-ab6131a999ca
	sigs.k8s.io/testing_frameworks => sigs.k8s.io/testing_frameworks v0.1.1
	sigs.k8s.io/yaml => sigs.k8s.io/yaml v1.1.0
)
