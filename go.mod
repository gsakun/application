module github.com/hd-Li/application

go 1.13

require (
	github.com/coreos/prometheus-operator v0.25.0 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/hd-Li/types v0.0.0-20200108072342-40227b4a545d
	github.com/knative/pkg v0.0.0-20190817231834-12ee58e32cc8
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/prometheus/common v0.9.1
	github.com/rancher/norman v0.0.0-20190319175355-e10534b012b0
	github.com/sirupsen/logrus v1.4.2
	github.com/snowzach/rotatefilehook v0.0.0-20180327172521-2f64f265f58c // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/apiserver v0.0.0-20181026185746-f1e867e1a455 // indirect
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/kubernetes v0.0.0-00010101000000-000000000000
)

replace (
	github.com/hd-Li/types => github.com/gsakun/types v0.0.0-20200429101124-60f0cee560f9
	github.com/knative/pkg => github.com/gsakun/pkg v0.0.0-20200421071615-21c5df62549f
	k8s.io/api => k8s.io/api v0.0.0-20181004124137-fd83cbc87e76
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20181004124836-1748dfb29e8a
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20180913025736-6dd46049f395
	k8s.io/client-go => k8s.io/client-go v9.0.0+incompatible
	k8s.io/kubernetes => k8s.io/kubernetes v1.12.2
)
