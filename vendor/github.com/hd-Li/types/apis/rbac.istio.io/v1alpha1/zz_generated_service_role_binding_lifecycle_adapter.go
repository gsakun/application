package v1alpha1

import (
	"github.com/hd-Li/types/pkg/istio/apis/rbac/v1alpha1"
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type ServiceRoleBindingLifecycle interface {
	Create(obj *v1alpha1.ServiceRoleBinding) (runtime.Object, error)
	Remove(obj *v1alpha1.ServiceRoleBinding) (runtime.Object, error)
	Updated(obj *v1alpha1.ServiceRoleBinding) (runtime.Object, error)
}

type serviceRoleBindingLifecycleAdapter struct {
	lifecycle ServiceRoleBindingLifecycle
}

func (w *serviceRoleBindingLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *serviceRoleBindingLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *serviceRoleBindingLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v1alpha1.ServiceRoleBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *serviceRoleBindingLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v1alpha1.ServiceRoleBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *serviceRoleBindingLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v1alpha1.ServiceRoleBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewServiceRoleBindingLifecycleAdapter(name string, clusterScoped bool, client ServiceRoleBindingInterface, l ServiceRoleBindingLifecycle) ServiceRoleBindingHandlerFunc {
	adapter := &serviceRoleBindingLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1alpha1.ServiceRoleBinding) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
