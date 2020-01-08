module github.com/hd-Li/application

go 1.13

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/coreos/prometheus-operator v0.25.0 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/hd-Li/types v0.0.0-00010101000000-000000000000
	github.com/knative/pkg v0.0.0-20190817231834-12ee58e32cc8
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/prometheus/client_model v0.1.0 // indirect
	github.com/prometheus/common v0.7.0 // indirect
	github.com/prometheus/procfs v0.0.8 // indirect
	github.com/rancher/norman v0.0.0-20190319175355-e10534b012b0
	k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
)

replace (
	github.com/hd-Li/types => ../types
	github.com/knative/pkg => github.com/rancher/pkg v0.0.0-20190514055449-b30ab9de040e
	k8s.io/api => k8s.io/api v0.0.0-20181004124137-fd83cbc87e76
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20181004124836-1748dfb29e8a
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20180913025736-6dd46049f395
	k8s.io/client-go => k8s.io/client-go v9.0.0+incompatible
)
