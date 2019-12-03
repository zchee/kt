module github.com/zchee/kt

go 1.13

require (
	cloud.google.com/go v0.45.1 // indirect
	github.com/cenkalti/backoff/v3 v3.1.1
	github.com/cornelk/hashmap v1.0.1
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-logr/logr v0.1.1-0.20190903151443-a1ebd699b195
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/google/go-cmp v0.3.2-0.20190829225427-b1c9c4891a65
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/minio/sha256-simd v0.1.2-0.20190917233721-f675151bb5e1
	github.com/panjf2000/ants/v2 v2.2.3-0.20191108040053-562ae1caf1f3
	github.com/prometheus/client_golang v1.2.1 // indirect
	github.com/spf13/cobra v0.0.6-0.20190805155617-b80588d523ec
	github.com/spf13/pflag v1.0.5
	github.com/zchee/color/v2 v2.0.3
	github.com/zeebo/xxh3 v0.0.0-20191021174148-b56a7dc3d80c
	go.opencensus.io v0.22.2
	go.uber.org/multierr v1.4.0
	go.uber.org/zap v1.12.0
	golang.org/x/crypto v0.0.0-20191108234033-bd318be0434a // indirect
	golang.org/x/net v0.0.0-20191109021931-daa7c04131f5 // indirect
	golang.org/x/sys v0.0.0-20191105231009-c1f44814a5cd // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/tools v0.0.0-20191108193012-7d206e10da11 // indirect
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/yaml.v2 v2.2.5 // indirect
	k8s.io/api v0.0.0-20191109101513-0171b7c15da1
	k8s.io/apiextensions-apiserver v0.0.0-20191109110701-3fdecfd8e730 // indirect
	k8s.io/apimachinery v0.0.0-20191109100838-fee41ff082ed
	k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
	sigs.k8s.io/controller-runtime v0.4.0
)

// pin
replace k8s.io/client-go => k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90 // kubernetes-1.16.0

replace (
	// k8s.io/client-go dependencies
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20181025213731-e84da0312774 // k8s.io/client-go@kubernetes-1.16.0
	golang.org/x/lint => golang.org/x/lint v0.0.0-20181217174547-8f45f776aaf1 // k8s.io/client-go@kubernetes-1.16.0
	golang.org/x/oauth2 => golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a // k8s.io/client-go@kubernetes-1.16.0
	golang.org/x/sync => golang.org/x/sync v0.0.0-20181108010431-42b317875d0f // k8s.io/client-go@kubernetes-1.16.0
	golang.org/x/sys => golang.org/x/sys v0.0.0-20190209173611-3b5209105503 // k8s.io/client-go@kubernetes-1.16.0
	golang.org/x/text => golang.org/x/text v0.3.1-0.20181227161524-e6919f6577db // k8s.io/client-go@kubernetes-1.16.0
	golang.org/x/time => golang.org/x/time v0.0.0-20161028155119-f51c12702a4d // k8s.io/client-go@kubernetes-1.16.0
	k8s.io/api => k8s.io/api v0.0.0-20190918155943-95b840bb6a1f // k8s.io/client-go@kubernetes-1.16.0

	// sigs.k8s.io/controller-runtime dependencies
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190918161926-8f644eb6e783 // sigs.k8s.io/controller-runtime@v0.4.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655 // k8s.io/client-go@kubernetes-1.16.0
	k8s.io/utils => k8s.io/utils v0.0.0-20190801114015-581e00157fb1 // k8s.io/client-go@kubernetes-1.16.0
)
