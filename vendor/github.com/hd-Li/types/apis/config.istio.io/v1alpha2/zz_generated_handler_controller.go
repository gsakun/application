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
	HandlerGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "Handler",
	}
	HandlerResource = metav1.APIResource{
		Name:         "handlers",
		SingularName: "handler",
		Namespaced:   true,

		Kind: HandlerGroupVersionKind.Kind,
	}
)

func NewHandler(namespace, name string, obj v1alpha2.Handler) *v1alpha2.Handler {
	obj.APIVersion, obj.Kind = HandlerGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type HandlerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []v1alpha2.Handler
}

type HandlerHandlerFunc func(key string, obj *v1alpha2.Handler) (runtime.Object, error)

type HandlerChangeHandlerFunc func(obj *v1alpha2.Handler) (runtime.Object, error)

type HandlerLister interface {
	List(namespace string, selector labels.Selector) (ret []*v1alpha2.Handler, err error)
	Get(namespace, name string) (*v1alpha2.Handler, error)
}

type HandlerController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() HandlerLister
	AddHandler(ctx context.Context, name string, handler HandlerHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler HandlerHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type HandlerInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*v1alpha2.Handler) (*v1alpha2.Handler, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha2.Handler, error)
	Get(name string, opts metav1.GetOptions) (*v1alpha2.Handler, error)
	Update(*v1alpha2.Handler) (*v1alpha2.Handler, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*HandlerList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() HandlerController
	AddHandler(ctx context.Context, name string, sync HandlerHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle HandlerLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync HandlerHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle HandlerLifecycle)
}

type handlerLister struct {
	controller *handlerController
}

func (l *handlerLister) List(namespace string, selector labels.Selector) (ret []*v1alpha2.Handler, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*v1alpha2.Handler))
	})
	return
}

func (l *handlerLister) Get(namespace, name string) (*v1alpha2.Handler, error) {
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
			Group:    HandlerGroupVersionKind.Group,
			Resource: "handler",
		}, key)
	}
	return obj.(*v1alpha2.Handler), nil
}

type handlerController struct {
	controller.GenericController
}

func (c *handlerController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *handlerController) Lister() HandlerLister {
	return &handlerLister{
		controller: c,
	}
}

func (c *handlerController) AddHandler(ctx context.Context, name string, handler HandlerHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha2.Handler); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *handlerController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler HandlerHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha2.Handler); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type handlerFactory struct {
}

func (c handlerFactory) Object() runtime.Object {
	return &v1alpha2.Handler{}
}

func (c handlerFactory) List() runtime.Object {
	return &HandlerList{}
}

func (s *handlerClient) Controller() HandlerController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.handlerControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(HandlerGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &handlerController{
		GenericController: genericController,
	}

	s.client.handlerControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type handlerClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   HandlerController
}

func (s *handlerClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *handlerClient) Create(o *v1alpha2.Handler) (*v1alpha2.Handler, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*v1alpha2.Handler), err
}

func (s *handlerClient) Get(name string, opts metav1.GetOptions) (*v1alpha2.Handler, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*v1alpha2.Handler), err
}

func (s *handlerClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha2.Handler, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*v1alpha2.Handler), err
}

func (s *handlerClient) Update(o *v1alpha2.Handler) (*v1alpha2.Handler, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*v1alpha2.Handler), err
}

func (s *handlerClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *handlerClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *handlerClient) List(opts metav1.ListOptions) (*HandlerList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*HandlerList), err
}

func (s *handlerClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *handlerClient) Patch(o *v1alpha2.Handler, patchType types.PatchType, data []byte, subresources ...string) (*v1alpha2.Handler, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*v1alpha2.Handler), err
}

func (s *handlerClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *handlerClient) AddHandler(ctx context.Context, name string, sync HandlerHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *handlerClient) AddLifecycle(ctx context.Context, name string, lifecycle HandlerLifecycle) {
	sync := NewHandlerLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *handlerClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync HandlerHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *handlerClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle HandlerLifecycle) {
	sync := NewHandlerLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type HandlerIndexer func(obj *v1alpha2.Handler) ([]string, error)

type HandlerClientCache interface {
	Get(namespace, name string) (*v1alpha2.Handler, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha2.Handler, error)

	Index(name string, indexer HandlerIndexer)
	GetIndexed(name, key string) ([]*v1alpha2.Handler, error)
}

type HandlerClient interface {
	Create(*v1alpha2.Handler) (*v1alpha2.Handler, error)
	Get(namespace, name string, opts metav1.GetOptions) (*v1alpha2.Handler, error)
	Update(*v1alpha2.Handler) (*v1alpha2.Handler, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*HandlerList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() HandlerClientCache

	OnCreate(ctx context.Context, name string, sync HandlerChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync HandlerChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync HandlerChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() HandlerInterface
}

type handlerClientCache struct {
	client *handlerClient2
}

type handlerClient2 struct {
	iface      HandlerInterface
	controller HandlerController
}

func (n *handlerClient2) Interface() HandlerInterface {
	return n.iface
}

func (n *handlerClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *handlerClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *handlerClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *handlerClient2) Create(obj *v1alpha2.Handler) (*v1alpha2.Handler, error) {
	return n.iface.Create(obj)
}

func (n *handlerClient2) Get(namespace, name string, opts metav1.GetOptions) (*v1alpha2.Handler, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *handlerClient2) Update(obj *v1alpha2.Handler) (*v1alpha2.Handler, error) {
	return n.iface.Update(obj)
}

func (n *handlerClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *handlerClient2) List(namespace string, opts metav1.ListOptions) (*HandlerList, error) {
	return n.iface.List(opts)
}

func (n *handlerClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *handlerClientCache) Get(namespace, name string) (*v1alpha2.Handler, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *handlerClientCache) List(namespace string, selector labels.Selector) ([]*v1alpha2.Handler, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *handlerClient2) Cache() HandlerClientCache {
	n.loadController()
	return &handlerClientCache{
		client: n,
	}
}

func (n *handlerClient2) OnCreate(ctx context.Context, name string, sync HandlerChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &handlerLifecycleDelegate{create: sync})
}

func (n *handlerClient2) OnChange(ctx context.Context, name string, sync HandlerChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &handlerLifecycleDelegate{update: sync})
}

func (n *handlerClient2) OnRemove(ctx context.Context, name string, sync HandlerChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &handlerLifecycleDelegate{remove: sync})
}

func (n *handlerClientCache) Index(name string, indexer HandlerIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*v1alpha2.Handler); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *handlerClientCache) GetIndexed(name, key string) ([]*v1alpha2.Handler, error) {
	var result []*v1alpha2.Handler
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*v1alpha2.Handler); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *handlerClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type handlerLifecycleDelegate struct {
	create HandlerChangeHandlerFunc
	update HandlerChangeHandlerFunc
	remove HandlerChangeHandlerFunc
}

func (n *handlerLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *handlerLifecycleDelegate) Create(obj *v1alpha2.Handler) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *handlerLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *handlerLifecycleDelegate) Remove(obj *v1alpha2.Handler) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *handlerLifecycleDelegate) Updated(obj *v1alpha2.Handler) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
