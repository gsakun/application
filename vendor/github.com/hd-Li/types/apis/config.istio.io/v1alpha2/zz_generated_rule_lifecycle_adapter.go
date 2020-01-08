package v1alpha2

import (
	"github.com/hd-Li/types/pkg/istio/apis/config/v1alpha2"
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type RuleLifecycle interface {
	Create(obj *v1alpha2.Rule) (runtime.Object, error)
	Remove(obj *v1alpha2.Rule) (runtime.Object, error)
	Updated(obj *v1alpha2.Rule) (runtime.Object, error)
}

type ruleLifecycleAdapter struct {
	lifecycle RuleLifecycle
}

func (w *ruleLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *ruleLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *ruleLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v1alpha2.Rule))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *ruleLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v1alpha2.Rule))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *ruleLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v1alpha2.Rule))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewRuleLifecycleAdapter(name string, clusterScoped bool, client RuleInterface, l RuleLifecycle) RuleHandlerFunc {
	adapter := &ruleLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1alpha2.Rule) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
