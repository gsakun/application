package v1alpha1

import (
	"github.com/hd-Li/types/pkg/istio/apis/rbac/v1alpha1"
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type ClusterRbacConfigLifecycle interface {
	Create(obj *v1alpha1.ClusterRbacConfig) (runtime.Object, error)
	Remove(obj *v1alpha1.ClusterRbacConfig) (runtime.Object, error)
	Updated(obj *v1alpha1.ClusterRbacConfig) (runtime.Object, error)
}

type clusterRbacConfigLifecycleAdapter struct {
	lifecycle ClusterRbacConfigLifecycle
}

func (w *clusterRbacConfigLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *clusterRbacConfigLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *clusterRbacConfigLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v1alpha1.ClusterRbacConfig))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterRbacConfigLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v1alpha1.ClusterRbacConfig))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterRbacConfigLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v1alpha1.ClusterRbacConfig))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewClusterRbacConfigLifecycleAdapter(name string, clusterScoped bool, client ClusterRbacConfigInterface, l ClusterRbacConfigLifecycle) ClusterRbacConfigHandlerFunc {
	adapter := &clusterRbacConfigLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1alpha1.ClusterRbacConfig) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
