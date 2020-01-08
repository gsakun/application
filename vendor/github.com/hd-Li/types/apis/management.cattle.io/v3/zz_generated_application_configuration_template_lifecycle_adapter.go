package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type ApplicationConfigurationTemplateLifecycle interface {
	Create(obj *ApplicationConfigurationTemplate) (runtime.Object, error)
	Remove(obj *ApplicationConfigurationTemplate) (runtime.Object, error)
	Updated(obj *ApplicationConfigurationTemplate) (runtime.Object, error)
}

type applicationConfigurationTemplateLifecycleAdapter struct {
	lifecycle ApplicationConfigurationTemplateLifecycle
}

func (w *applicationConfigurationTemplateLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *applicationConfigurationTemplateLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *applicationConfigurationTemplateLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*ApplicationConfigurationTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *applicationConfigurationTemplateLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*ApplicationConfigurationTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *applicationConfigurationTemplateLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*ApplicationConfigurationTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewApplicationConfigurationTemplateLifecycleAdapter(name string, clusterScoped bool, client ApplicationConfigurationTemplateInterface, l ApplicationConfigurationTemplateLifecycle) ApplicationConfigurationTemplateHandlerFunc {
	adapter := &applicationConfigurationTemplateLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *ApplicationConfigurationTemplate) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
