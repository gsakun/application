package v1alpha1

import (
	"github.com/hd-Li/types/pkg/istio/apis/rbac/v1alpha1"
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type ServiceRoleLifecycle interface {
	Create(obj *v1alpha1.ServiceRole) (runtime.Object, error)
	Remove(obj *v1alpha1.ServiceRole) (runtime.Object, error)
	Updated(obj *v1alpha1.ServiceRole) (runtime.Object, error)
}

type serviceRoleLifecycleAdapter struct {
	lifecycle ServiceRoleLifecycle
}

func (w *serviceRoleLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *serviceRoleLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *serviceRoleLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v1alpha1.ServiceRole))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *serviceRoleLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v1alpha1.ServiceRole))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *serviceRoleLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v1alpha1.ServiceRole))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewServiceRoleLifecycleAdapter(name string, clusterScoped bool, client ServiceRoleInterface, l ServiceRoleLifecycle) ServiceRoleHandlerFunc {
	adapter := &serviceRoleLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1alpha1.ServiceRole) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
