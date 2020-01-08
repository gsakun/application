package v1alpha2

import (
	"github.com/hd-Li/types/pkg/istio/apis/config/v1alpha2"
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type QuotaSpecLifecycle interface {
	Create(obj *v1alpha2.QuotaSpec) (runtime.Object, error)
	Remove(obj *v1alpha2.QuotaSpec) (runtime.Object, error)
	Updated(obj *v1alpha2.QuotaSpec) (runtime.Object, error)
}

type quotaSpecLifecycleAdapter struct {
	lifecycle QuotaSpecLifecycle
}

func (w *quotaSpecLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *quotaSpecLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *quotaSpecLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v1alpha2.QuotaSpec))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *quotaSpecLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v1alpha2.QuotaSpec))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *quotaSpecLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v1alpha2.QuotaSpec))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewQuotaSpecLifecycleAdapter(name string, clusterScoped bool, client QuotaSpecInterface, l QuotaSpecLifecycle) QuotaSpecHandlerFunc {
	adapter := &quotaSpecLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1alpha2.QuotaSpec) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
