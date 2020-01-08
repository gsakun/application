package v1alpha2

import (
	"context"

	"github.com/hd-Li/types/pkg/istio/apis/config/v1alpha2"
	"github.com/rancher/norman/controller"
	"github.com/rancher/norman/objectclient"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	RuleGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "Rule",
	}
	RuleResource = metav1.APIResource{
		Name:         "rules",
		SingularName: "rule",
		Namespaced:   true,

		Kind: RuleGroupVersionKind.Kind,
	}
)

func NewRule(namespace, name string, obj v1alpha2.Rule) *v1alpha2.Rule {
	obj.APIVersion, obj.Kind = RuleGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type RuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []v1alpha2.Rule
}

type RuleHandlerFunc func(key string, obj *v1alpha2.Rule) (runtime.Object, error)

type RuleChangeHandlerFunc func(obj *v1alpha2.Rule) (runtime.Object, error)

type RuleLister interface {
	List(namespace string, selector labels.Selector) (ret []*v1alpha2.Rule, err error)
	Get(namespace, name string) (*v1alpha2.Rule, error)
}

type RuleController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() RuleLister
	AddHandler(ctx context.Context, name string, handler RuleHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler RuleHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type RuleInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*v1alpha2.Rule) (*v1alpha2.Rule, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha2.Rule, error)
	Get(name string, opts metav1.GetOptions) (*v1alpha2.Rule, error)
	Update(*v1alpha2.Rule) (*v1alpha2.Rule, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*RuleList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() RuleController
	AddHandler(ctx context.Context, name string, sync RuleHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle RuleLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync RuleHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle RuleLifecycle)
}

type ruleLister struct {
	controller *ruleController
}

func (l *ruleLister) List(namespace string, selector labels.Selector) (ret []*v1alpha2.Rule, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*v1alpha2.Rule))
	})
	return
}

func (l *ruleLister) Get(namespace, name string) (*v1alpha2.Rule, error) {
	var key string
	if namespace != "" {
		key = namespace + "/" + name
	} else {
		key = name
	}
	obj, exists, err := l.controller.Informer().GetIndexer().GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schema.GroupResource{
			Group:    RuleGroupVersionKind.Group,
			Resource: "rule",
		}, key)
	}
	return obj.(*v1alpha2.Rule), nil
}

type ruleController struct {
	controller.GenericController
}

func (c *ruleController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *ruleController) Lister() RuleLister {
	return &ruleLister{
		controller: c,
	}
}

func (c *ruleController) AddHandler(ctx context.Context, name string, handler RuleHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha2.Rule); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *ruleController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler RuleHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha2.Rule); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type ruleFactory struct {
}

func (c ruleFactory) Object() runtime.Object {
	return &v1alpha2.Rule{}
}

func (c ruleFactory) List() runtime.Object {
	return &RuleList{}
}

func (s *ruleClient) Controller() RuleController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.ruleControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(RuleGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &ruleController{
		GenericController: genericController,
	}

	s.client.ruleControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type ruleClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   RuleController
}

func (s *ruleClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *ruleClient) Create(o *v1alpha2.Rule) (*v1alpha2.Rule, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*v1alpha2.Rule), err
}

func (s *ruleClient) Get(name string, opts metav1.GetOptions) (*v1alpha2.Rule, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*v1alpha2.Rule), err
}

func (s *ruleClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha2.Rule, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*v1alpha2.Rule), err
}

func (s *ruleClient) Update(o *v1alpha2.Rule) (*v1alpha2.Rule, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*v1alpha2.Rule), err
}

func (s *ruleClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *ruleClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *ruleClient) List(opts metav1.ListOptions) (*RuleList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*RuleList), err
}

func (s *ruleClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *ruleClient) Patch(o *v1alpha2.Rule, patchType types.PatchType, data []byte, subresources ...string) (*v1alpha2.Rule, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*v1alpha2.Rule), err
}

func (s *ruleClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *ruleClient) AddHandler(ctx context.Context, name string, sync RuleHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *ruleClient) AddLifecycle(ctx context.Context, name string, lifecycle RuleLifecycle) {
	sync := NewRuleLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *ruleClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync RuleHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *ruleClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle RuleLifecycle) {
	sync := NewRuleLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type RuleIndexer func(obj *v1alpha2.Rule) ([]string, error)

type RuleClientCache interface {
	Get(namespace, name string) (*v1alpha2.Rule, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha2.Rule, error)

	Index(name string, indexer RuleIndexer)
	GetIndexed(name, key string) ([]*v1alpha2.Rule, error)
}

type RuleClient interface {
	Create(*v1alpha2.Rule) (*v1alpha2.Rule, error)
	Get(namespace, name string, opts metav1.GetOptions) (*v1alpha2.Rule, error)
	Update(*v1alpha2.Rule) (*v1alpha2.Rule, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*RuleList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() RuleClientCache

	OnCreate(ctx context.Context, name string, sync RuleChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync RuleChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync RuleChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() RuleInterface
}

type ruleClientCache struct {
	client *ruleClient2
}

type ruleClient2 struct {
	iface      RuleInterface
	controller RuleController
}

func (n *ruleClient2) Interface() RuleInterface {
	return n.iface
}

func (n *ruleClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *ruleClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *ruleClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *ruleClient2) Create(obj *v1alpha2.Rule) (*v1alpha2.Rule, error) {
	return n.iface.Create(obj)
}

func (n *ruleClient2) Get(namespace, name string, opts metav1.GetOptions) (*v1alpha2.Rule, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *ruleClient2) Update(obj *v1alpha2.Rule) (*v1alpha2.Rule, error) {
	return n.iface.Update(obj)
}

func (n *ruleClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *ruleClient2) List(namespace string, opts metav1.ListOptions) (*RuleList, error) {
	return n.iface.List(opts)
}

func (n *ruleClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *ruleClientCache) Get(namespace, name string) (*v1alpha2.Rule, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *ruleClientCache) List(namespace string, selector labels.Selector) ([]*v1alpha2.Rule, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *ruleClient2) Cache() RuleClientCache {
	n.loadController()
	return &ruleClientCache{
		client: n,
	}
}

func (n *ruleClient2) OnCreate(ctx context.Context, name string, sync RuleChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &ruleLifecycleDelegate{create: sync})
}

func (n *ruleClient2) OnChange(ctx context.Context, name string, sync RuleChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &ruleLifecycleDelegate{update: sync})
}

func (n *ruleClient2) OnRemove(ctx context.Context, name string, sync RuleChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &ruleLifecycleDelegate{remove: sync})
}

func (n *ruleClientCache) Index(name string, indexer RuleIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*v1alpha2.Rule); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *ruleClientCache) GetIndexed(name, key string) ([]*v1alpha2.Rule, error) {
	var result []*v1alpha2.Rule
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*v1alpha2.Rule); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *ruleClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type ruleLifecycleDelegate struct {
	create RuleChangeHandlerFunc
	update RuleChangeHandlerFunc
	remove RuleChangeHandlerFunc
}

func (n *ruleLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *ruleLifecycleDelegate) Create(obj *v1alpha2.Rule) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *ruleLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *ruleLifecycleDelegate) Remove(obj *v1alpha2.Rule) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *ruleLifecycleDelegate) Updated(obj *v1alpha2.Rule) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
