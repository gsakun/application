package v1alpha2

import (
	"github.com/hd-Li/types/pkg/istio/apis/config/v1alpha2"
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type QuotaSpecBindingLifecycle interface {
	Create(obj *v1alpha2.QuotaSpecBinding) (runtime.Object, error)
	Remove(obj *v1alpha2.QuotaSpecBinding) (runtime.Object, error)
	Updated(obj *v1alpha2.QuotaSpecBinding) (runtime.Object, error)
}

type quotaSpecBindingLifecycleAdapter struct {
	lifecycle QuotaSpecBindingLifecycle
}

func (w *quotaSpecBindingLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *quotaSpecBindingLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *quotaSpecBindingLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v1alpha2.QuotaSpecBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *quotaSpecBindingLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v1alpha2.QuotaSpecBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *quotaSpecBindingLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v1alpha2.QuotaSpecBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewQuotaSpecBindingLifecycleAdapter(name string, clusterScoped bool, client QuotaSpecBindingInterface, l QuotaSpecBindingLifecycle) QuotaSpecBindingHandlerFunc {
	adapter := &quotaSpecBindingLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1alpha2.QuotaSpecBinding) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
