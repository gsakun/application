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
	ServiceRoleGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "ServiceRole",
	}
	ServiceRoleResource = metav1.APIResource{
		Name:         "serviceroles",
		SingularName: "servicerole",
		Namespaced:   true,

		Kind: ServiceRoleGroupVersionKind.Kind,
	}
)

func NewServiceRole(namespace, name string, obj v1alpha1.ServiceRole) *v1alpha1.ServiceRole {
	obj.APIVersion, obj.Kind = ServiceRoleGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type ServiceRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []v1alpha1.ServiceRole
}

type ServiceRoleHandlerFunc func(key string, obj *v1alpha1.ServiceRole) (runtime.Object, error)

type ServiceRoleChangeHandlerFunc func(obj *v1alpha1.ServiceRole) (runtime.Object, error)

type ServiceRoleLister interface {
	List(namespace string, selector labels.Selector) (ret []*v1alpha1.ServiceRole, err error)
	Get(namespace, name string) (*v1alpha1.ServiceRole, error)
}

type ServiceRoleController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() ServiceRoleLister
	AddHandler(ctx context.Context, name string, handler ServiceRoleHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler ServiceRoleHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type ServiceRoleInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*v1alpha1.ServiceRole) (*v1alpha1.ServiceRole, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha1.ServiceRole, error)
	Get(name string, opts metav1.GetOptions) (*v1alpha1.ServiceRole, error)
	Update(*v1alpha1.ServiceRole) (*v1alpha1.ServiceRole, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*ServiceRoleList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() ServiceRoleController
	AddHandler(ctx context.Context, name string, sync ServiceRoleHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle ServiceRoleLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ServiceRoleHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ServiceRoleLifecycle)
}

type serviceRoleLister struct {
	controller *serviceRoleController
}

func (l *serviceRoleLister) List(namespace string, selector labels.Selector) (ret []*v1alpha1.ServiceRole, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*v1alpha1.ServiceRole))
	})
	return
}

func (l *serviceRoleLister) Get(namespace, name string) (*v1alpha1.ServiceRole, error) {
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
			Group:    ServiceRoleGroupVersionKind.Group,
			Resource: "serviceRole",
		}, key)
	}
	return obj.(*v1alpha1.ServiceRole), nil
}

type serviceRoleController struct {
	controller.GenericController
}

func (c *serviceRoleController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *serviceRoleController) Lister() ServiceRoleLister {
	return &serviceRoleLister{
		controller: c,
	}
}

func (c *serviceRoleController) AddHandler(ctx context.Context, name string, handler ServiceRoleHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha1.ServiceRole); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *serviceRoleController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler ServiceRoleHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha1.ServiceRole); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type serviceRoleFactory struct {
}

func (c serviceRoleFactory) Object() runtime.Object {
	return &v1alpha1.ServiceRole{}
}

func (c serviceRoleFactory) List() runtime.Object {
	return &ServiceRoleList{}
}

func (s *serviceRoleClient) Controller() ServiceRoleController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.serviceRoleControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(ServiceRoleGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &serviceRoleController{
		GenericController: genericController,
	}

	s.client.serviceRoleControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type serviceRoleClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   ServiceRoleController
}

func (s *serviceRoleClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *serviceRoleClient) Create(o *v1alpha1.ServiceRole) (*v1alpha1.ServiceRole, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*v1alpha1.ServiceRole), err
}

func (s *serviceRoleClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.ServiceRole, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*v1alpha1.ServiceRole), err
}

func (s *serviceRoleClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha1.ServiceRole, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*v1alpha1.ServiceRole), err
}

func (s *serviceRoleClient) Update(o *v1alpha1.ServiceRole) (*v1alpha1.ServiceRole, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*v1alpha1.ServiceRole), err
}

func (s *serviceRoleClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *serviceRoleClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *serviceRoleClient) List(opts metav1.ListOptions) (*ServiceRoleList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*ServiceRoleList), err
}

func (s *serviceRoleClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *serviceRoleClient) Patch(o *v1alpha1.ServiceRole, patchType types.PatchType, data []byte, subresources ...string) (*v1alpha1.ServiceRole, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*v1alpha1.ServiceRole), err
}

func (s *serviceRoleClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *serviceRoleClient) AddHandler(ctx context.Context, name string, sync ServiceRoleHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *serviceRoleClient) AddLifecycle(ctx context.Context, name string, lifecycle ServiceRoleLifecycle) {
	sync := NewServiceRoleLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *serviceRoleClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ServiceRoleHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *serviceRoleClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ServiceRoleLifecycle) {
	sync := NewServiceRoleLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type ServiceRoleIndexer func(obj *v1alpha1.ServiceRole) ([]string, error)

type ServiceRoleClientCache interface {
	Get(namespace, name string) (*v1alpha1.ServiceRole, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha1.ServiceRole, error)

	Index(name string, indexer ServiceRoleIndexer)
	GetIndexed(name, key string) ([]*v1alpha1.ServiceRole, error)
}

type ServiceRoleClient interface {
	Create(*v1alpha1.ServiceRole) (*v1alpha1.ServiceRole, error)
	Get(namespace, name string, opts metav1.GetOptions) (*v1alpha1.ServiceRole, error)
	Update(*v1alpha1.ServiceRole) (*v1alpha1.ServiceRole, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*ServiceRoleList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() ServiceRoleClientCache

	OnCreate(ctx context.Context, name string, sync ServiceRoleChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync ServiceRoleChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync ServiceRoleChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() ServiceRoleInterface
}

type serviceRoleClientCache struct {
	client *serviceRoleClient2
}

type serviceRoleClient2 struct {
	iface      ServiceRoleInterface
	controller ServiceRoleController
}

func (n *serviceRoleClient2) Interface() ServiceRoleInterface {
	return n.iface
}

func (n *serviceRoleClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *serviceRoleClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *serviceRoleClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *serviceRoleClient2) Create(obj *v1alpha1.ServiceRole) (*v1alpha1.ServiceRole, error) {
	return n.iface.Create(obj)
}

func (n *serviceRoleClient2) Get(namespace, name string, opts metav1.GetOptions) (*v1alpha1.ServiceRole, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *serviceRoleClient2) Update(obj *v1alpha1.ServiceRole) (*v1alpha1.ServiceRole, error) {
	return n.iface.Update(obj)
}

func (n *serviceRoleClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *serviceRoleClient2) List(namespace string, opts metav1.ListOptions) (*ServiceRoleList, error) {
	return n.iface.List(opts)
}

func (n *serviceRoleClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *serviceRoleClientCache) Get(namespace, name string) (*v1alpha1.ServiceRole, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *serviceRoleClientCache) List(namespace string, selector labels.Selector) ([]*v1alpha1.ServiceRole, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *serviceRoleClient2) Cache() ServiceRoleClientCache {
	n.loadController()
	return &serviceRoleClientCache{
		client: n,
	}
}

func (n *serviceRoleClient2) OnCreate(ctx context.Context, name string, sync ServiceRoleChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &serviceRoleLifecycleDelegate{create: sync})
}

func (n *serviceRoleClient2) OnChange(ctx context.Context, name string, sync ServiceRoleChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &serviceRoleLifecycleDelegate{update: sync})
}

func (n *serviceRoleClient2) OnRemove(ctx context.Context, name string, sync ServiceRoleChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &serviceRoleLifecycleDelegate{remove: sync})
}

func (n *serviceRoleClientCache) Index(name string, indexer ServiceRoleIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*v1alpha1.ServiceRole); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *serviceRoleClientCache) GetIndexed(name, key string) ([]*v1alpha1.ServiceRole, error) {
	var result []*v1alpha1.ServiceRole
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*v1alpha1.ServiceRole); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *serviceRoleClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type serviceRoleLifecycleDelegate struct {
	create ServiceRoleChangeHandlerFunc
	update ServiceRoleChangeHandlerFunc
	remove ServiceRoleChangeHandlerFunc
}

func (n *serviceRoleLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *serviceRoleLifecycleDelegate) Create(obj *v1alpha1.ServiceRole) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *serviceRoleLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *serviceRoleLifecycleDelegate) Remove(obj *v1alpha1.ServiceRole) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *serviceRoleLifecycleDelegate) Updated(obj *v1alpha1.ServiceRole) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
