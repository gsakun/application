package v3

import (
	"context"

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
	ApplicationGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "Application",
	}
	ApplicationResource = metav1.APIResource{
		Name:         "applications",
		SingularName: "application",
		Namespaced:   true,

		Kind: ApplicationGroupVersionKind.Kind,
	}
)

func NewApplication(namespace, name string, obj Application) *Application {
	obj.APIVersion, obj.Kind = ApplicationGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application
}

type ApplicationHandlerFunc func(key string, obj *Application) (runtime.Object, error)

type ApplicationChangeHandlerFunc func(obj *Application) (runtime.Object, error)

type ApplicationLister interface {
	List(namespace string, selector labels.Selector) (ret []*Application, err error)
	Get(namespace, name string) (*Application, error)
}

type ApplicationController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() ApplicationLister
	AddHandler(ctx context.Context, name string, handler ApplicationHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler ApplicationHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type ApplicationInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*Application) (*Application, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*Application, error)
	Get(name string, opts metav1.GetOptions) (*Application, error)
	Update(*Application) (*Application, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*ApplicationList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() ApplicationController
	AddHandler(ctx context.Context, name string, sync ApplicationHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle ApplicationLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ApplicationHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ApplicationLifecycle)
}

type applicationLister struct {
	controller *applicationController
}

func (l *applicationLister) List(namespace string, selector labels.Selector) (ret []*Application, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*Application))
	})
	return
}

func (l *applicationLister) Get(namespace, name string) (*Application, error) {
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
			Group:    ApplicationGroupVersionKind.Group,
			Resource: "application",
		}, key)
	}
	return obj.(*Application), nil
}

type applicationController struct {
	controller.GenericController
}

func (c *applicationController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *applicationController) Lister() ApplicationLister {
	return &applicationLister{
		controller: c,
	}
}

func (c *applicationController) AddHandler(ctx context.Context, name string, handler ApplicationHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Application); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *applicationController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler ApplicationHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*Application); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type applicationFactory struct {
}

func (c applicationFactory) Object() runtime.Object {
	return &Application{}
}

func (c applicationFactory) List() runtime.Object {
	return &ApplicationList{}
}

func (s *applicationClient) Controller() ApplicationController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.applicationControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(ApplicationGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &applicationController{
		GenericController: genericController,
	}

	s.client.applicationControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type applicationClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   ApplicationController
}

func (s *applicationClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *applicationClient) Create(o *Application) (*Application, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*Application), err
}

func (s *applicationClient) Get(name string, opts metav1.GetOptions) (*Application, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*Application), err
}

func (s *applicationClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*Application, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*Application), err
}

func (s *applicationClient) Update(o *Application) (*Application, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*Application), err
}

func (s *applicationClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *applicationClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *applicationClient) List(opts metav1.ListOptions) (*ApplicationList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*ApplicationList), err
}

func (s *applicationClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *applicationClient) Patch(o *Application, patchType types.PatchType, data []byte, subresources ...string) (*Application, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*Application), err
}

func (s *applicationClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *applicationClient) AddHandler(ctx context.Context, name string, sync ApplicationHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *applicationClient) AddLifecycle(ctx context.Context, name string, lifecycle ApplicationLifecycle) {
	sync := NewApplicationLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *applicationClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ApplicationHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *applicationClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ApplicationLifecycle) {
	sync := NewApplicationLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type ApplicationIndexer func(obj *Application) ([]string, error)

type ApplicationClientCache interface {
	Get(namespace, name string) (*Application, error)
	List(namespace string, selector labels.Selector) ([]*Application, error)

	Index(name string, indexer ApplicationIndexer)
	GetIndexed(name, key string) ([]*Application, error)
}

type ApplicationClient interface {
	Create(*Application) (*Application, error)
	Get(namespace, name string, opts metav1.GetOptions) (*Application, error)
	Update(*Application) (*Application, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*ApplicationList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() ApplicationClientCache

	OnCreate(ctx context.Context, name string, sync ApplicationChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync ApplicationChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync ApplicationChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() ApplicationInterface
}

type applicationClientCache struct {
	client *applicationClient2
}

type applicationClient2 struct {
	iface      ApplicationInterface
	controller ApplicationController
}

func (n *applicationClient2) Interface() ApplicationInterface {
	return n.iface
}

func (n *applicationClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *applicationClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *applicationClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *applicationClient2) Create(obj *Application) (*Application, error) {
	return n.iface.Create(obj)
}

func (n *applicationClient2) Get(namespace, name string, opts metav1.GetOptions) (*Application, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *applicationClient2) Update(obj *Application) (*Application, error) {
	return n.iface.Update(obj)
}

func (n *applicationClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *applicationClient2) List(namespace string, opts metav1.ListOptions) (*ApplicationList, error) {
	return n.iface.List(opts)
}

func (n *applicationClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *applicationClientCache) Get(namespace, name string) (*Application, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *applicationClientCache) List(namespace string, selector labels.Selector) ([]*Application, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *applicationClient2) Cache() ApplicationClientCache {
	n.loadController()
	return &applicationClientCache{
		client: n,
	}
}

func (n *applicationClient2) OnCreate(ctx context.Context, name string, sync ApplicationChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &applicationLifecycleDelegate{create: sync})
}

func (n *applicationClient2) OnChange(ctx context.Context, name string, sync ApplicationChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &applicationLifecycleDelegate{update: sync})
}

func (n *applicationClient2) OnRemove(ctx context.Context, name string, sync ApplicationChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &applicationLifecycleDelegate{remove: sync})
}

func (n *applicationClientCache) Index(name string, indexer ApplicationIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*Application); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *applicationClientCache) GetIndexed(name, key string) ([]*Application, error) {
	var result []*Application
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*Application); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *applicationClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type applicationLifecycleDelegate struct {
	create ApplicationChangeHandlerFunc
	update ApplicationChangeHandlerFunc
	remove ApplicationChangeHandlerFunc
}

func (n *applicationLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *applicationLifecycleDelegate) Create(obj *Application) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *applicationLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *applicationLifecycleDelegate) Remove(obj *Application) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *applicationLifecycleDelegate) Updated(obj *Application) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
