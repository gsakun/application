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
	QuotaSpecBindingGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "QuotaSpecBinding",
	}
	QuotaSpecBindingResource = metav1.APIResource{
		Name:         "quotaspecbindings",
		SingularName: "quotaspecbinding",
		Namespaced:   true,

		Kind: QuotaSpecBindingGroupVersionKind.Kind,
	}
)

func NewQuotaSpecBinding(namespace, name string, obj v1alpha2.QuotaSpecBinding) *v1alpha2.QuotaSpecBinding {
	obj.APIVersion, obj.Kind = QuotaSpecBindingGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type QuotaSpecBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []v1alpha2.QuotaSpecBinding
}

type QuotaSpecBindingHandlerFunc func(key string, obj *v1alpha2.QuotaSpecBinding) (runtime.Object, error)

type QuotaSpecBindingChangeHandlerFunc func(obj *v1alpha2.QuotaSpecBinding) (runtime.Object, error)

type QuotaSpecBindingLister interface {
	List(namespace string, selector labels.Selector) (ret []*v1alpha2.QuotaSpecBinding, err error)
	Get(namespace, name string) (*v1alpha2.QuotaSpecBinding, error)
}

type QuotaSpecBindingController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() QuotaSpecBindingLister
	AddHandler(ctx context.Context, name string, handler QuotaSpecBindingHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler QuotaSpecBindingHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type QuotaSpecBindingInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*v1alpha2.QuotaSpecBinding) (*v1alpha2.QuotaSpecBinding, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha2.QuotaSpecBinding, error)
	Get(name string, opts metav1.GetOptions) (*v1alpha2.QuotaSpecBinding, error)
	Update(*v1alpha2.QuotaSpecBinding) (*v1alpha2.QuotaSpecBinding, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*QuotaSpecBindingList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() QuotaSpecBindingController
	AddHandler(ctx context.Context, name string, sync QuotaSpecBindingHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle QuotaSpecBindingLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync QuotaSpecBindingHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle QuotaSpecBindingLifecycle)
}

type quotaSpecBindingLister struct {
	controller *quotaSpecBindingController
}

func (l *quotaSpecBindingLister) List(namespace string, selector labels.Selector) (ret []*v1alpha2.QuotaSpecBinding, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*v1alpha2.QuotaSpecBinding))
	})
	return
}

func (l *quotaSpecBindingLister) Get(namespace, name string) (*v1alpha2.QuotaSpecBinding, error) {
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
			Group:    QuotaSpecBindingGroupVersionKind.Group,
			Resource: "quotaSpecBinding",
		}, key)
	}
	return obj.(*v1alpha2.QuotaSpecBinding), nil
}

type quotaSpecBindingController struct {
	controller.GenericController
}

func (c *quotaSpecBindingController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *quotaSpecBindingController) Lister() QuotaSpecBindingLister {
	return &quotaSpecBindingLister{
		controller: c,
	}
}

func (c *quotaSpecBindingController) AddHandler(ctx context.Context, name string, handler QuotaSpecBindingHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha2.QuotaSpecBinding); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *quotaSpecBindingController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler QuotaSpecBindingHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha2.QuotaSpecBinding); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type quotaSpecBindingFactory struct {
}

func (c quotaSpecBindingFactory) Object() runtime.Object {
	return &v1alpha2.QuotaSpecBinding{}
}

func (c quotaSpecBindingFactory) List() runtime.Object {
	return &QuotaSpecBindingList{}
}

func (s *quotaSpecBindingClient) Controller() QuotaSpecBindingController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.quotaSpecBindingControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(QuotaSpecBindingGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &quotaSpecBindingController{
		GenericController: genericController,
	}

	s.client.quotaSpecBindingControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type quotaSpecBindingClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   QuotaSpecBindingController
}

func (s *quotaSpecBindingClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *quotaSpecBindingClient) Create(o *v1alpha2.QuotaSpecBinding) (*v1alpha2.QuotaSpecBinding, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*v1alpha2.QuotaSpecBinding), err
}

func (s *quotaSpecBindingClient) Get(name string, opts metav1.GetOptions) (*v1alpha2.QuotaSpecBinding, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*v1alpha2.QuotaSpecBinding), err
}

func (s *quotaSpecBindingClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha2.QuotaSpecBinding, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*v1alpha2.QuotaSpecBinding), err
}

func (s *quotaSpecBindingClient) Update(o *v1alpha2.QuotaSpecBinding) (*v1alpha2.QuotaSpecBinding, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*v1alpha2.QuotaSpecBinding), err
}

func (s *quotaSpecBindingClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *quotaSpecBindingClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *quotaSpecBindingClient) List(opts metav1.ListOptions) (*QuotaSpecBindingList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*QuotaSpecBindingList), err
}

func (s *quotaSpecBindingClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *quotaSpecBindingClient) Patch(o *v1alpha2.QuotaSpecBinding, patchType types.PatchType, data []byte, subresources ...string) (*v1alpha2.QuotaSpecBinding, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*v1alpha2.QuotaSpecBinding), err
}

func (s *quotaSpecBindingClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *quotaSpecBindingClient) AddHandler(ctx context.Context, name string, sync QuotaSpecBindingHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *quotaSpecBindingClient) AddLifecycle(ctx context.Context, name string, lifecycle QuotaSpecBindingLifecycle) {
	sync := NewQuotaSpecBindingLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *quotaSpecBindingClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync QuotaSpecBindingHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *quotaSpecBindingClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle QuotaSpecBindingLifecycle) {
	sync := NewQuotaSpecBindingLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type QuotaSpecBindingIndexer func(obj *v1alpha2.QuotaSpecBinding) ([]string, error)

type QuotaSpecBindingClientCache interface {
	Get(namespace, name string) (*v1alpha2.QuotaSpecBinding, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha2.QuotaSpecBinding, error)

	Index(name string, indexer QuotaSpecBindingIndexer)
	GetIndexed(name, key string) ([]*v1alpha2.QuotaSpecBinding, error)
}

type QuotaSpecBindingClient interface {
	Create(*v1alpha2.QuotaSpecBinding) (*v1alpha2.QuotaSpecBinding, error)
	Get(namespace, name string, opts metav1.GetOptions) (*v1alpha2.QuotaSpecBinding, error)
	Update(*v1alpha2.QuotaSpecBinding) (*v1alpha2.QuotaSpecBinding, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*QuotaSpecBindingList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() QuotaSpecBindingClientCache

	OnCreate(ctx context.Context, name string, sync QuotaSpecBindingChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync QuotaSpecBindingChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync QuotaSpecBindingChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() QuotaSpecBindingInterface
}

type quotaSpecBindingClientCache struct {
	client *quotaSpecBindingClient2
}

type quotaSpecBindingClient2 struct {
	iface      QuotaSpecBindingInterface
	controller QuotaSpecBindingController
}

func (n *quotaSpecBindingClient2) Interface() QuotaSpecBindingInterface {
	return n.iface
}

func (n *quotaSpecBindingClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *quotaSpecBindingClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *quotaSpecBindingClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *quotaSpecBindingClient2) Create(obj *v1alpha2.QuotaSpecBinding) (*v1alpha2.QuotaSpecBinding, error) {
	return n.iface.Create(obj)
}

func (n *quotaSpecBindingClient2) Get(namespace, name string, opts metav1.GetOptions) (*v1alpha2.QuotaSpecBinding, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *quotaSpecBindingClient2) Update(obj *v1alpha2.QuotaSpecBinding) (*v1alpha2.QuotaSpecBinding, error) {
	return n.iface.Update(obj)
}

func (n *quotaSpecBindingClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *quotaSpecBindingClient2) List(namespace string, opts metav1.ListOptions) (*QuotaSpecBindingList, error) {
	return n.iface.List(opts)
}

func (n *quotaSpecBindingClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *quotaSpecBindingClientCache) Get(namespace, name string) (*v1alpha2.QuotaSpecBinding, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *quotaSpecBindingClientCache) List(namespace string, selector labels.Selector) ([]*v1alpha2.QuotaSpecBinding, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *quotaSpecBindingClient2) Cache() QuotaSpecBindingClientCache {
	n.loadController()
	return &quotaSpecBindingClientCache{
		client: n,
	}
}

func (n *quotaSpecBindingClient2) OnCreate(ctx context.Context, name string, sync QuotaSpecBindingChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &quotaSpecBindingLifecycleDelegate{create: sync})
}

func (n *quotaSpecBindingClient2) OnChange(ctx context.Context, name string, sync QuotaSpecBindingChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &quotaSpecBindingLifecycleDelegate{update: sync})
}

func (n *quotaSpecBindingClient2) OnRemove(ctx context.Context, name string, sync QuotaSpecBindingChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &quotaSpecBindingLifecycleDelegate{remove: sync})
}

func (n *quotaSpecBindingClientCache) Index(name string, indexer QuotaSpecBindingIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*v1alpha2.QuotaSpecBinding); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *quotaSpecBindingClientCache) GetIndexed(name, key string) ([]*v1alpha2.QuotaSpecBinding, error) {
	var result []*v1alpha2.QuotaSpecBinding
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*v1alpha2.QuotaSpecBinding); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *quotaSpecBindingClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type quotaSpecBindingLifecycleDelegate struct {
	create QuotaSpecBindingChangeHandlerFunc
	update QuotaSpecBindingChangeHandlerFunc
	remove QuotaSpecBindingChangeHandlerFunc
}

func (n *quotaSpecBindingLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *quotaSpecBindingLifecycleDelegate) Create(obj *v1alpha2.QuotaSpecBinding) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *quotaSpecBindingLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *quotaSpecBindingLifecycleDelegate) Remove(obj *v1alpha2.QuotaSpecBinding) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *quotaSpecBindingLifecycleDelegate) Updated(obj *v1alpha2.QuotaSpecBinding) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
