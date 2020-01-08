package v1alpha2

import (
	"github.com/hd-Li/types/pkg/istio/apis/config/v1alpha2"
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type InstanceLifecycle interface {
	Create(obj *v1alpha2.Instance) (runtime.Object, error)
	Remove(obj *v1alpha2.Instance) (runtime.Object, error)
	Updated(obj *v1alpha2.Instance) (runtime.Object, error)
}

type instanceLifecycleAdapter struct {
	lifecycle InstanceLifecycle
}

func (w *instanceLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *instanceLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *instanceLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v1alpha2.Instance))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *instanceLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v1alpha2.Instance))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *instanceLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v1alpha2.Instance))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewInstanceLifecycleAdapter(name string, clusterScoped bool, client InstanceInterface, l InstanceLifecycle) InstanceHandlerFunc {
	adapter := &instanceLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1alpha2.Instance) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
