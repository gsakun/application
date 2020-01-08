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
	QuotaSpecGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "QuotaSpec",
	}
	QuotaSpecResource = metav1.APIResource{
		Name:         "quotaspecs",
		SingularName: "quotaspec",
		Namespaced:   true,

		Kind: QuotaSpecGroupVersionKind.Kind,
	}
)

func NewQuotaSpec(namespace, name string, obj v1alpha2.QuotaSpec) *v1alpha2.QuotaSpec {
	obj.APIVersion, obj.Kind = QuotaSpecGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type QuotaSpecList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []v1alpha2.QuotaSpec
}

type QuotaSpecHandlerFunc func(key string, obj *v1alpha2.QuotaSpec) (runtime.Object, error)

type QuotaSpecChangeHandlerFunc func(obj *v1alpha2.QuotaSpec) (runtime.Object, error)

type QuotaSpecLister interface {
	List(namespace string, selector labels.Selector) (ret []*v1alpha2.QuotaSpec, err error)
	Get(namespace, name string) (*v1alpha2.QuotaSpec, error)
}

type QuotaSpecController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() QuotaSpecLister
	AddHandler(ctx context.Context, name string, handler QuotaSpecHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler QuotaSpecHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type QuotaSpecInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*v1alpha2.QuotaSpec) (*v1alpha2.QuotaSpec, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha2.QuotaSpec, error)
	Get(name string, opts metav1.GetOptions) (*v1alpha2.QuotaSpec, error)
	Update(*v1alpha2.QuotaSpec) (*v1alpha2.QuotaSpec, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*QuotaSpecList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() QuotaSpecController
	AddHandler(ctx context.Context, name string, sync QuotaSpecHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle QuotaSpecLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync QuotaSpecHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle QuotaSpecLifecycle)
}

type quotaSpecLister struct {
	controller *quotaSpecController
}

func (l *quotaSpecLister) List(namespace string, selector labels.Selector) (ret []*v1alpha2.QuotaSpec, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*v1alpha2.QuotaSpec))
	})
	return
}

func (l *quotaSpecLister) Get(namespace, name string) (*v1alpha2.QuotaSpec, error) {
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
			Group:    QuotaSpecGroupVersionKind.Group,
			Resource: "quotaSpec",
		}, key)
	}
	return obj.(*v1alpha2.QuotaSpec), nil
}

type quotaSpecController struct {
	controller.GenericController
}

func (c *quotaSpecController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *quotaSpecController) Lister() QuotaSpecLister {
	return &quotaSpecLister{
		controller: c,
	}
}

func (c *quotaSpecController) AddHandler(ctx context.Context, name string, handler QuotaSpecHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha2.QuotaSpec); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *quotaSpecController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler QuotaSpecHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha2.QuotaSpec); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type quotaSpecFactory struct {
}

func (c quotaSpecFactory) Object() runtime.Object {
	return &v1alpha2.QuotaSpec{}
}

func (c quotaSpecFactory) List() runtime.Object {
	return &QuotaSpecList{}
}

func (s *quotaSpecClient) Controller() QuotaSpecController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.quotaSpecControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(QuotaSpecGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &quotaSpecController{
		GenericController: genericController,
	}

	s.client.quotaSpecControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type quotaSpecClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   QuotaSpecController
}

func (s *quotaSpecClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *quotaSpecClient) Create(o *v1alpha2.QuotaSpec) (*v1alpha2.QuotaSpec, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*v1alpha2.QuotaSpec), err
}

func (s *quotaSpecClient) Get(name string, opts metav1.GetOptions) (*v1alpha2.QuotaSpec, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*v1alpha2.QuotaSpec), err
}

func (s *quotaSpecClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha2.QuotaSpec, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*v1alpha2.QuotaSpec), err
}

func (s *quotaSpecClient) Update(o *v1alpha2.QuotaSpec) (*v1alpha2.QuotaSpec, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*v1alpha2.QuotaSpec), err
}

func (s *quotaSpecClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *quotaSpecClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *quotaSpecClient) List(opts metav1.ListOptions) (*QuotaSpecList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*QuotaSpecList), err
}

func (s *quotaSpecClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *quotaSpecClient) Patch(o *v1alpha2.QuotaSpec, patchType types.PatchType, data []byte, subresources ...string) (*v1alpha2.QuotaSpec, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*v1alpha2.QuotaSpec), err
}

func (s *quotaSpecClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *quotaSpecClient) AddHandler(ctx context.Context, name string, sync QuotaSpecHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *quotaSpecClient) AddLifecycle(ctx context.Context, name string, lifecycle QuotaSpecLifecycle) {
	sync := NewQuotaSpecLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *quotaSpecClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync QuotaSpecHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *quotaSpecClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle QuotaSpecLifecycle) {
	sync := NewQuotaSpecLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type QuotaSpecIndexer func(obj *v1alpha2.QuotaSpec) ([]string, error)

type QuotaSpecClientCache interface {
	Get(namespace, name string) (*v1alpha2.QuotaSpec, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha2.QuotaSpec, error)

	Index(name string, indexer QuotaSpecIndexer)
	GetIndexed(name, key string) ([]*v1alpha2.QuotaSpec, error)
}

type QuotaSpecClient interface {
	Create(*v1alpha2.QuotaSpec) (*v1alpha2.QuotaSpec, error)
	Get(namespace, name string, opts metav1.GetOptions) (*v1alpha2.QuotaSpec, error)
	Update(*v1alpha2.QuotaSpec) (*v1alpha2.QuotaSpec, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*QuotaSpecList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() QuotaSpecClientCache

	OnCreate(ctx context.Context, name string, sync QuotaSpecChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync QuotaSpecChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync QuotaSpecChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() QuotaSpecInterface
}

type quotaSpecClientCache struct {
	client *quotaSpecClient2
}

type quotaSpecClient2 struct {
	iface      QuotaSpecInterface
	controller QuotaSpecController
}

func (n *quotaSpecClient2) Interface() QuotaSpecInterface {
	return n.iface
}

func (n *quotaSpecClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *quotaSpecClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *quotaSpecClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *quotaSpecClient2) Create(obj *v1alpha2.QuotaSpec) (*v1alpha2.QuotaSpec, error) {
	return n.iface.Create(obj)
}

func (n *quotaSpecClient2) Get(namespace, name string, opts metav1.GetOptions) (*v1alpha2.QuotaSpec, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *quotaSpecClient2) Update(obj *v1alpha2.QuotaSpec) (*v1alpha2.QuotaSpec, error) {
	return n.iface.Update(obj)
}

func (n *quotaSpecClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *quotaSpecClient2) List(namespace string, opts metav1.ListOptions) (*QuotaSpecList, error) {
	return n.iface.List(opts)
}

func (n *quotaSpecClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *quotaSpecClientCache) Get(namespace, name string) (*v1alpha2.QuotaSpec, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *quotaSpecClientCache) List(namespace string, selector labels.Selector) ([]*v1alpha2.QuotaSpec, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *quotaSpecClient2) Cache() QuotaSpecClientCache {
	n.loadController()
	return &quotaSpecClientCache{
		client: n,
	}
}

func (n *quotaSpecClient2) OnCreate(ctx context.Context, name string, sync QuotaSpecChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &quotaSpecLifecycleDelegate{create: sync})
}

func (n *quotaSpecClient2) OnChange(ctx context.Context, name string, sync QuotaSpecChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &quotaSpecLifecycleDelegate{update: sync})
}

func (n *quotaSpecClient2) OnRemove(ctx context.Context, name string, sync QuotaSpecChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &quotaSpecLifecycleDelegate{remove: sync})
}

func (n *quotaSpecClientCache) Index(name string, indexer QuotaSpecIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*v1alpha2.QuotaSpec); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *quotaSpecClientCache) GetIndexed(name, key string) ([]*v1alpha2.QuotaSpec, error) {
	var result []*v1alpha2.QuotaSpec
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*v1alpha2.QuotaSpec); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *quotaSpecClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type quotaSpecLifecycleDelegate struct {
	create QuotaSpecChangeHandlerFunc
	update QuotaSpecChangeHandlerFunc
	remove QuotaSpecChangeHandlerFunc
}

func (n *quotaSpecLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *quotaSpecLifecycleDelegate) Create(obj *v1alpha2.QuotaSpec) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *quotaSpecLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *quotaSpecLifecycleDelegate) Remove(obj *v1alpha2.QuotaSpec) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *quotaSpecLifecycleDelegate) Updated(obj *v1alpha2.QuotaSpec) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
