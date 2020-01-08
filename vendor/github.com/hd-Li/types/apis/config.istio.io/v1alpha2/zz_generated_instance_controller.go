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
	InstanceGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "Instance",
	}
	InstanceResource = metav1.APIResource{
		Name:         "instances",
		SingularName: "instance",
		Namespaced:   true,

		Kind: InstanceGroupVersionKind.Kind,
	}
)

func NewInstance(namespace, name string, obj v1alpha2.Instance) *v1alpha2.Instance {
	obj.APIVersion, obj.Kind = InstanceGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []v1alpha2.Instance
}

type InstanceHandlerFunc func(key string, obj *v1alpha2.Instance) (runtime.Object, error)

type InstanceChangeHandlerFunc func(obj *v1alpha2.Instance) (runtime.Object, error)

type InstanceLister interface {
	List(namespace string, selector labels.Selector) (ret []*v1alpha2.Instance, err error)
	Get(namespace, name string) (*v1alpha2.Instance, error)
}

type InstanceController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() InstanceLister
	AddHandler(ctx context.Context, name string, handler InstanceHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler InstanceHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type InstanceInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*v1alpha2.Instance) (*v1alpha2.Instance, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha2.Instance, error)
	Get(name string, opts metav1.GetOptions) (*v1alpha2.Instance, error)
	Update(*v1alpha2.Instance) (*v1alpha2.Instance, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*InstanceList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() InstanceController
	AddHandler(ctx context.Context, name string, sync InstanceHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle InstanceLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync InstanceHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle InstanceLifecycle)
}

type instanceLister struct {
	controller *instanceController
}

func (l *instanceLister) List(namespace string, selector labels.Selector) (ret []*v1alpha2.Instance, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*v1alpha2.Instance))
	})
	return
}

func (l *instanceLister) Get(namespace, name string) (*v1alpha2.Instance, error) {
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
			Group:    InstanceGroupVersionKind.Group,
			Resource: "instance",
		}, key)
	}
	return obj.(*v1alpha2.Instance), nil
}

type instanceController struct {
	controller.GenericController
}

func (c *instanceController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *instanceController) Lister() InstanceLister {
	return &instanceLister{
		controller: c,
	}
}

func (c *instanceController) AddHandler(ctx context.Context, name string, handler InstanceHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha2.Instance); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *instanceController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler InstanceHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha2.Instance); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type instanceFactory struct {
}

func (c instanceFactory) Object() runtime.Object {
	return &v1alpha2.Instance{}
}

func (c instanceFactory) List() runtime.Object {
	return &InstanceList{}
}

func (s *instanceClient) Controller() InstanceController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.instanceControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(InstanceGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &instanceController{
		GenericController: genericController,
	}

	s.client.instanceControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type instanceClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   InstanceController
}

func (s *instanceClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *instanceClient) Create(o *v1alpha2.Instance) (*v1alpha2.Instance, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*v1alpha2.Instance), err
}

func (s *instanceClient) Get(name string, opts metav1.GetOptions) (*v1alpha2.Instance, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*v1alpha2.Instance), err
}

func (s *instanceClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha2.Instance, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*v1alpha2.Instance), err
}

func (s *instanceClient) Update(o *v1alpha2.Instance) (*v1alpha2.Instance, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*v1alpha2.Instance), err
}

func (s *instanceClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *instanceClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *instanceClient) List(opts metav1.ListOptions) (*InstanceList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*InstanceList), err
}

func (s *instanceClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *instanceClient) Patch(o *v1alpha2.Instance, patchType types.PatchType, data []byte, subresources ...string) (*v1alpha2.Instance, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*v1alpha2.Instance), err
}

func (s *instanceClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *instanceClient) AddHandler(ctx context.Context, name string, sync InstanceHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *instanceClient) AddLifecycle(ctx context.Context, name string, lifecycle InstanceLifecycle) {
	sync := NewInstanceLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *instanceClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync InstanceHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *instanceClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle InstanceLifecycle) {
	sync := NewInstanceLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type InstanceIndexer func(obj *v1alpha2.Instance) ([]string, error)

type InstanceClientCache interface {
	Get(namespace, name string) (*v1alpha2.Instance, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha2.Instance, error)

	Index(name string, indexer InstanceIndexer)
	GetIndexed(name, key string) ([]*v1alpha2.Instance, error)
}

type InstanceClient interface {
	Create(*v1alpha2.Instance) (*v1alpha2.Instance, error)
	Get(namespace, name string, opts metav1.GetOptions) (*v1alpha2.Instance, error)
	Update(*v1alpha2.Instance) (*v1alpha2.Instance, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*InstanceList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() InstanceClientCache

	OnCreate(ctx context.Context, name string, sync InstanceChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync InstanceChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync InstanceChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() InstanceInterface
}

type instanceClientCache struct {
	client *instanceClient2
}

type instanceClient2 struct {
	iface      InstanceInterface
	controller InstanceController
}

func (n *instanceClient2) Interface() InstanceInterface {
	return n.iface
}

func (n *instanceClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *instanceClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *instanceClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *instanceClient2) Create(obj *v1alpha2.Instance) (*v1alpha2.Instance, error) {
	return n.iface.Create(obj)
}

func (n *instanceClient2) Get(namespace, name string, opts metav1.GetOptions) (*v1alpha2.Instance, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *instanceClient2) Update(obj *v1alpha2.Instance) (*v1alpha2.Instance, error) {
	return n.iface.Update(obj)
}

func (n *instanceClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *instanceClient2) List(namespace string, opts metav1.ListOptions) (*InstanceList, error) {
	return n.iface.List(opts)
}

func (n *instanceClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *instanceClientCache) Get(namespace, name string) (*v1alpha2.Instance, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *instanceClientCache) List(namespace string, selector labels.Selector) ([]*v1alpha2.Instance, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *instanceClient2) Cache() InstanceClientCache {
	n.loadController()
	return &instanceClientCache{
		client: n,
	}
}

func (n *instanceClient2) OnCreate(ctx context.Context, name string, sync InstanceChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &instanceLifecycleDelegate{create: sync})
}

func (n *instanceClient2) OnChange(ctx context.Context, name string, sync InstanceChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &instanceLifecycleDelegate{update: sync})
}

func (n *instanceClient2) OnRemove(ctx context.Context, name string, sync InstanceChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &instanceLifecycleDelegate{remove: sync})
}

func (n *instanceClientCache) Index(name string, indexer InstanceIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*v1alpha2.Instance); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *instanceClientCache) GetIndexed(name, key string) ([]*v1alpha2.Instance, error) {
	var result []*v1alpha2.Instance
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*v1alpha2.Instance); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *instanceClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type instanceLifecycleDelegate struct {
	create InstanceChangeHandlerFunc
	update InstanceChangeHandlerFunc
	remove InstanceChangeHandlerFunc
}

func (n *instanceLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *instanceLifecycleDelegate) Create(obj *v1alpha2.Instance) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *instanceLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *instanceLifecycleDelegate) Remove(obj *v1alpha2.Instance) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *instanceLifecycleDelegate) Updated(obj *v1alpha2.Instance) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
