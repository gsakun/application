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
	ClusterRbacConfigGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "ClusterRbacConfig",
	}
	ClusterRbacConfigResource = metav1.APIResource{
		Name:         "clusterrbacconfigs",
		SingularName: "clusterrbacconfig",
		Namespaced:   true,

		Kind: ClusterRbacConfigGroupVersionKind.Kind,
	}
)

func NewClusterRbacConfig(namespace, name string, obj v1alpha1.ClusterRbacConfig) *v1alpha1.ClusterRbacConfig {
	obj.APIVersion, obj.Kind = ClusterRbacConfigGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type ClusterRbacConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []v1alpha1.ClusterRbacConfig
}

type ClusterRbacConfigHandlerFunc func(key string, obj *v1alpha1.ClusterRbacConfig) (runtime.Object, error)

type ClusterRbacConfigChangeHandlerFunc func(obj *v1alpha1.ClusterRbacConfig) (runtime.Object, error)

type ClusterRbacConfigLister interface {
	List(namespace string, selector labels.Selector) (ret []*v1alpha1.ClusterRbacConfig, err error)
	Get(namespace, name string) (*v1alpha1.ClusterRbacConfig, error)
}

type ClusterRbacConfigController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() ClusterRbacConfigLister
	AddHandler(ctx context.Context, name string, handler ClusterRbacConfigHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler ClusterRbacConfigHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type ClusterRbacConfigInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*v1alpha1.ClusterRbacConfig) (*v1alpha1.ClusterRbacConfig, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha1.ClusterRbacConfig, error)
	Get(name string, opts metav1.GetOptions) (*v1alpha1.ClusterRbacConfig, error)
	Update(*v1alpha1.ClusterRbacConfig) (*v1alpha1.ClusterRbacConfig, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*ClusterRbacConfigList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() ClusterRbacConfigController
	AddHandler(ctx context.Context, name string, sync ClusterRbacConfigHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle ClusterRbacConfigLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ClusterRbacConfigHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ClusterRbacConfigLifecycle)
}

type clusterRbacConfigLister struct {
	controller *clusterRbacConfigController
}

func (l *clusterRbacConfigLister) List(namespace string, selector labels.Selector) (ret []*v1alpha1.ClusterRbacConfig, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*v1alpha1.ClusterRbacConfig))
	})
	return
}

func (l *clusterRbacConfigLister) Get(namespace, name string) (*v1alpha1.ClusterRbacConfig, error) {
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
			Group:    ClusterRbacConfigGroupVersionKind.Group,
			Resource: "clusterRbacConfig",
		}, key)
	}
	return obj.(*v1alpha1.ClusterRbacConfig), nil
}

type clusterRbacConfigController struct {
	controller.GenericController
}

func (c *clusterRbacConfigController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *clusterRbacConfigController) Lister() ClusterRbacConfigLister {
	return &clusterRbacConfigLister{
		controller: c,
	}
}

func (c *clusterRbacConfigController) AddHandler(ctx context.Context, name string, handler ClusterRbacConfigHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha1.ClusterRbacConfig); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *clusterRbacConfigController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler ClusterRbacConfigHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha1.ClusterRbacConfig); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type clusterRbacConfigFactory struct {
}

func (c clusterRbacConfigFactory) Object() runtime.Object {
	return &v1alpha1.ClusterRbacConfig{}
}

func (c clusterRbacConfigFactory) List() runtime.Object {
	return &ClusterRbacConfigList{}
}

func (s *clusterRbacConfigClient) Controller() ClusterRbacConfigController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.clusterRbacConfigControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(ClusterRbacConfigGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &clusterRbacConfigController{
		GenericController: genericController,
	}

	s.client.clusterRbacConfigControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type clusterRbacConfigClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   ClusterRbacConfigController
}

func (s *clusterRbacConfigClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *clusterRbacConfigClient) Create(o *v1alpha1.ClusterRbacConfig) (*v1alpha1.ClusterRbacConfig, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*v1alpha1.ClusterRbacConfig), err
}

func (s *clusterRbacConfigClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.ClusterRbacConfig, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*v1alpha1.ClusterRbacConfig), err
}

func (s *clusterRbacConfigClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha1.ClusterRbacConfig, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*v1alpha1.ClusterRbacConfig), err
}

func (s *clusterRbacConfigClient) Update(o *v1alpha1.ClusterRbacConfig) (*v1alpha1.ClusterRbacConfig, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*v1alpha1.ClusterRbacConfig), err
}

func (s *clusterRbacConfigClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *clusterRbacConfigClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *clusterRbacConfigClient) List(opts metav1.ListOptions) (*ClusterRbacConfigList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*ClusterRbacConfigList), err
}

func (s *clusterRbacConfigClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *clusterRbacConfigClient) Patch(o *v1alpha1.ClusterRbacConfig, patchType types.PatchType, data []byte, subresources ...string) (*v1alpha1.ClusterRbacConfig, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*v1alpha1.ClusterRbacConfig), err
}

func (s *clusterRbacConfigClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *clusterRbacConfigClient) AddHandler(ctx context.Context, name string, sync ClusterRbacConfigHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *clusterRbacConfigClient) AddLifecycle(ctx context.Context, name string, lifecycle ClusterRbacConfigLifecycle) {
	sync := NewClusterRbacConfigLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *clusterRbacConfigClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ClusterRbacConfigHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *clusterRbacConfigClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ClusterRbacConfigLifecycle) {
	sync := NewClusterRbacConfigLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type ClusterRbacConfigIndexer func(obj *v1alpha1.ClusterRbacConfig) ([]string, error)

type ClusterRbacConfigClientCache interface {
	Get(namespace, name string) (*v1alpha1.ClusterRbacConfig, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha1.ClusterRbacConfig, error)

	Index(name string, indexer ClusterRbacConfigIndexer)
	GetIndexed(name, key string) ([]*v1alpha1.ClusterRbacConfig, error)
}

type ClusterRbacConfigClient interface {
	Create(*v1alpha1.ClusterRbacConfig) (*v1alpha1.ClusterRbacConfig, error)
	Get(namespace, name string, opts metav1.GetOptions) (*v1alpha1.ClusterRbacConfig, error)
	Update(*v1alpha1.ClusterRbacConfig) (*v1alpha1.ClusterRbacConfig, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*ClusterRbacConfigList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() ClusterRbacConfigClientCache

	OnCreate(ctx context.Context, name string, sync ClusterRbacConfigChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync ClusterRbacConfigChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync ClusterRbacConfigChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() ClusterRbacConfigInterface
}

type clusterRbacConfigClientCache struct {
	client *clusterRbacConfigClient2
}

type clusterRbacConfigClient2 struct {
	iface      ClusterRbacConfigInterface
	controller ClusterRbacConfigController
}

func (n *clusterRbacConfigClient2) Interface() ClusterRbacConfigInterface {
	return n.iface
}

func (n *clusterRbacConfigClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *clusterRbacConfigClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *clusterRbacConfigClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *clusterRbacConfigClient2) Create(obj *v1alpha1.ClusterRbacConfig) (*v1alpha1.ClusterRbacConfig, error) {
	return n.iface.Create(obj)
}

func (n *clusterRbacConfigClient2) Get(namespace, name string, opts metav1.GetOptions) (*v1alpha1.ClusterRbacConfig, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *clusterRbacConfigClient2) Update(obj *v1alpha1.ClusterRbacConfig) (*v1alpha1.ClusterRbacConfig, error) {
	return n.iface.Update(obj)
}

func (n *clusterRbacConfigClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *clusterRbacConfigClient2) List(namespace string, opts metav1.ListOptions) (*ClusterRbacConfigList, error) {
	return n.iface.List(opts)
}

func (n *clusterRbacConfigClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *clusterRbacConfigClientCache) Get(namespace, name string) (*v1alpha1.ClusterRbacConfig, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *clusterRbacConfigClientCache) List(namespace string, selector labels.Selector) ([]*v1alpha1.ClusterRbacConfig, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *clusterRbacConfigClient2) Cache() ClusterRbacConfigClientCache {
	n.loadController()
	return &clusterRbacConfigClientCache{
		client: n,
	}
}

func (n *clusterRbacConfigClient2) OnCreate(ctx context.Context, name string, sync ClusterRbacConfigChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &clusterRbacConfigLifecycleDelegate{create: sync})
}

func (n *clusterRbacConfigClient2) OnChange(ctx context.Context, name string, sync ClusterRbacConfigChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &clusterRbacConfigLifecycleDelegate{update: sync})
}

func (n *clusterRbacConfigClient2) OnRemove(ctx context.Context, name string, sync ClusterRbacConfigChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &clusterRbacConfigLifecycleDelegate{remove: sync})
}

func (n *clusterRbacConfigClientCache) Index(name string, indexer ClusterRbacConfigIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*v1alpha1.ClusterRbacConfig); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *clusterRbacConfigClientCache) GetIndexed(name, key string) ([]*v1alpha1.ClusterRbacConfig, error) {
	var result []*v1alpha1.ClusterRbacConfig
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*v1alpha1.ClusterRbacConfig); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *clusterRbacConfigClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type clusterRbacConfigLifecycleDelegate struct {
	create ClusterRbacConfigChangeHandlerFunc
	update ClusterRbacConfigChangeHandlerFunc
	remove ClusterRbacConfigChangeHandlerFunc
}

func (n *clusterRbacConfigLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *clusterRbacConfigLifecycleDelegate) Create(obj *v1alpha1.ClusterRbacConfig) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *clusterRbacConfigLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *clusterRbacConfigLifecycleDelegate) Remove(obj *v1alpha1.ClusterRbacConfig) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *clusterRbacConfigLifecycleDelegate) Updated(obj *v1alpha1.ClusterRbacConfig) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
