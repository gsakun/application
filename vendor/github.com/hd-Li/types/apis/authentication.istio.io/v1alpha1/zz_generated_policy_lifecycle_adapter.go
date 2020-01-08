package v1alpha1

import (
	"github.com/hd-Li/types/pkg/istio/apis/authentication/v1alpha1"
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type PolicyLifecycle interface {
	Create(obj *v1alpha1.Policy) (runtime.Object, error)
	Remove(obj *v1alpha1.Policy) (runtime.Object, error)
	Updated(obj *v1alpha1.Policy) (runtime.Object, error)
}

type policyLifecycleAdapter struct {
	lifecycle PolicyLifecycle
}

func (w *policyLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *policyLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *policyLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v1alpha1.Policy))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *policyLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v1alpha1.Policy))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *policyLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v1alpha1.Policy))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewPolicyLifecycleAdapter(name string, clusterScoped bool, client PolicyInterface, l PolicyLifecycle) PolicyHandlerFunc {
	adapter := &policyLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1alpha1.Policy) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
