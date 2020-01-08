package v1alpha1

import (
	"context"

	"github.com/hd-Li/types/pkg/istio/apis/rbac/v1alpha1"
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
	ServiceRoleBindingGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "ServiceRoleBinding",
	}
	ServiceRoleBindingResource = metav1.APIResource{
		Name:         "servicerolebindings",
		SingularName: "servicerolebinding",
		Namespaced:   true,

		Kind: ServiceRoleBindingGroupVersionKind.Kind,
	}
)

func NewServiceRoleBinding(namespace, name string, obj v1alpha1.ServiceRoleBinding) *v1alpha1.ServiceRoleBinding {
	obj.APIVersion, obj.Kind = ServiceRoleBindingGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type ServiceRoleBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []v1alpha1.ServiceRoleBinding
}

type ServiceRoleBindingHandlerFunc func(key string, obj *v1alpha1.ServiceRoleBinding) (runtime.Object, error)

type ServiceRoleBindingChangeHandlerFunc func(obj *v1alpha1.ServiceRoleBinding) (runtime.Object, error)

type ServiceRoleBindingLister interface {
	List(namespace string, selector labels.Selector) (ret []*v1alpha1.ServiceRoleBinding, err error)
	Get(namespace, name string) (*v1alpha1.ServiceRoleBinding, error)
}

type ServiceRoleBindingController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() ServiceRoleBindingLister
	AddHandler(ctx context.Context, name string, handler ServiceRoleBindingHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler ServiceRoleBindingHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type ServiceRoleBindingInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*v1alpha1.ServiceRoleBinding) (*v1alpha1.ServiceRoleBinding, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha1.ServiceRoleBinding, error)
	Get(name string, opts metav1.GetOptions) (*v1alpha1.ServiceRoleBinding, error)
	Update(*v1alpha1.ServiceRoleBinding) (*v1alpha1.ServiceRoleBinding, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*ServiceRoleBindingList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() ServiceRoleBindingController
	AddHandler(ctx context.Context, name string, sync ServiceRoleBindingHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle ServiceRoleBindingLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ServiceRoleBindingHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ServiceRoleBindingLifecycle)
}

type serviceRoleBindingLister struct {
	controller *serviceRoleBindingController
}

func (l *serviceRoleBindingLister) List(namespace string, selector labels.Selector) (ret []*v1alpha1.ServiceRoleBinding, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*v1alpha1.ServiceRoleBinding))
	})
	return
}

func (l *serviceRoleBindingLister) Get(namespace, name string) (*v1alpha1.ServiceRoleBinding, error) {
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
			Group:    ServiceRoleBindingGroupVersionKind.Group,
			Resource: "serviceRoleBinding",
		}, key)
	}
	return obj.(*v1alpha1.ServiceRoleBinding), nil
}

type serviceRoleBindingController struct {
	controller.GenericController
}

func (c *serviceRoleBindingController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *serviceRoleBindingController) Lister() ServiceRoleBindingLister {
	return &serviceRoleBindingLister{
		controller: c,
	}
}

func (c *serviceRoleBindingController) AddHandler(ctx context.Context, name string, handler ServiceRoleBindingHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha1.ServiceRoleBinding); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *serviceRoleBindingController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler ServiceRoleBindingHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha1.ServiceRoleBinding); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type serviceRoleBindingFactory struct {
}

func (c serviceRoleBindingFactory) Object() runtime.Object {
	return &v1alpha1.ServiceRoleBinding{}
}

func (c serviceRoleBindingFactory) List() runtime.Object {
	return &ServiceRoleBindingList{}
}

func (s *serviceRoleBindingClient) Controller() ServiceRoleBindingController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.serviceRoleBindingControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(ServiceRoleBindingGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &serviceRoleBindingController{
		GenericController: genericController,
	}

	s.client.serviceRoleBindingControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type serviceRoleBindingClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   ServiceRoleBindingController
}

func (s *serviceRoleBindingClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *serviceRoleBindingClient) Create(o *v1alpha1.ServiceRoleBinding) (*v1alpha1.ServiceRoleBinding, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*v1alpha1.ServiceRoleBinding), err
}

func (s *serviceRoleBindingClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.ServiceRoleBinding, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*v1alpha1.ServiceRoleBinding), err
}

func (s *serviceRoleBindingClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha1.ServiceRoleBinding, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*v1alpha1.ServiceRoleBinding), err
}

func (s *serviceRoleBindingClient) Update(o *v1alpha1.ServiceRoleBinding) (*v1alpha1.ServiceRoleBinding, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*v1alpha1.ServiceRoleBinding), err
}

func (s *serviceRoleBindingClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *serviceRoleBindingClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *serviceRoleBindingClient) List(opts metav1.ListOptions) (*ServiceRoleBindingList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*ServiceRoleBindingList), err
}

func (s *serviceRoleBindingClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *serviceRoleBindingClient) Patch(o *v1alpha1.ServiceRoleBinding, patchType types.PatchType, data []byte, subresources ...string) (*v1alpha1.ServiceRoleBinding, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*v1alpha1.ServiceRoleBinding), err
}

func (s *serviceRoleBindingClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *serviceRoleBindingClient) AddHandler(ctx context.Context, name string, sync ServiceRoleBindingHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *serviceRoleBindingClient) AddLifecycle(ctx context.Context, name string, lifecycle ServiceRoleBindingLifecycle) {
	sync := NewServiceRoleBindingLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *serviceRoleBindingClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ServiceRoleBindingHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *serviceRoleBindingClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ServiceRoleBindingLifecycle) {
	sync := NewServiceRoleBindingLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type ServiceRoleBindingIndexer func(obj *v1alpha1.ServiceRoleBinding) ([]string, error)

type ServiceRoleBindingClientCache interface {
	Get(namespace, name string) (*v1alpha1.ServiceRoleBinding, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha1.ServiceRoleBinding, error)

	Index(name string, indexer ServiceRoleBindingIndexer)
	GetIndexed(name, key string) ([]*v1alpha1.ServiceRoleBinding, error)
}

type ServiceRoleBindingClient interface {
	Create(*v1alpha1.ServiceRoleBinding) (*v1alpha1.ServiceRoleBinding, error)
	Get(namespace, name string, opts metav1.GetOptions) (*v1alpha1.ServiceRoleBinding, error)
	Update(*v1alpha1.ServiceRoleBinding) (*v1alpha1.ServiceRoleBinding, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*ServiceRoleBindingList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() ServiceRoleBindingClientCache

	OnCreate(ctx context.Context, name string, sync ServiceRoleBindingChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync ServiceRoleBindingChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync ServiceRoleBindingChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() ServiceRoleBindingInterface
}

type serviceRoleBindingClientCache struct {
	client *serviceRoleBindingClient2
}

type serviceRoleBindingClient2 struct {
	iface      ServiceRoleBindingInterface
	controller ServiceRoleBindingController
}

func (n *serviceRoleBindingClient2) Interface() ServiceRoleBindingInterface {
	return n.iface
}

func (n *serviceRoleBindingClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *serviceRoleBindingClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *serviceRoleBindingClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *serviceRoleBindingClient2) Create(obj *v1alpha1.ServiceRoleBinding) (*v1alpha1.ServiceRoleBinding, error) {
	return n.iface.Create(obj)
}

func (n *serviceRoleBindingClient2) Get(namespace, name string, opts metav1.GetOptions) (*v1alpha1.ServiceRoleBinding, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *serviceRoleBindingClient2) Update(obj *v1alpha1.ServiceRoleBinding) (*v1alpha1.ServiceRoleBinding, error) {
	return n.iface.Update(obj)
}

func (n *serviceRoleBindingClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *serviceRoleBindingClient2) List(namespace string, opts metav1.ListOptions) (*ServiceRoleBindingList, error) {
	return n.iface.List(opts)
}

func (n *serviceRoleBindingClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *serviceRoleBindingClientCache) Get(namespace, name string) (*v1alpha1.ServiceRoleBinding, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *serviceRoleBindingClientCache) List(namespace string, selector labels.Selector) ([]*v1alpha1.ServiceRoleBinding, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *serviceRoleBindingClient2) Cache() ServiceRoleBindingClientCache {
	n.loadController()
	return &serviceRoleBindingClientCache{
		client: n,
	}
}

func (n *serviceRoleBindingClient2) OnCreate(ctx context.Context, name string, sync ServiceRoleBindingChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &serviceRoleBindingLifecycleDelegate{create: sync})
}

func (n *serviceRoleBindingClient2) OnChange(ctx context.Context, name string, sync ServiceRoleBindingChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &serviceRoleBindingLifecycleDelegate{update: sync})
}

func (n *serviceRoleBindingClient2) OnRemove(ctx context.Context, name string, sync ServiceRoleBindingChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &serviceRoleBindingLifecycleDelegate{remove: sync})
}

func (n *serviceRoleBindingClientCache) Index(name string, indexer ServiceRoleBindingIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*v1alpha1.ServiceRoleBinding); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *serviceRoleBindingClientCache) GetIndexed(name, key string) ([]*v1alpha1.ServiceRoleBinding, error) {
	var result []*v1alpha1.ServiceRoleBinding
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*v1alpha1.ServiceRoleBinding); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *serviceRoleBindingClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type serviceRoleBindingLifecycleDelegate struct {
	create ServiceRoleBindingChangeHandlerFunc
	update ServiceRoleBindingChangeHandlerFunc
	remove ServiceRoleBindingChangeHandlerFunc
}

func (n *serviceRoleBindingLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *serviceRoleBindingLifecycleDelegate) Create(obj *v1alpha1.ServiceRoleBinding) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *serviceRoleBindingLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *serviceRoleBindingLifecycleDelegate) Remove(obj *v1alpha1.ServiceRoleBinding) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *serviceRoleBindingLifecycleDelegate) Updated(obj *v1alpha1.ServiceRoleBinding) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
