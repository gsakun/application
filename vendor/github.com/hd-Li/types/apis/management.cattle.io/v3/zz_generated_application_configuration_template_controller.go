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
	ApplicationConfigurationTemplateGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "ApplicationConfigurationTemplate",
	}
	ApplicationConfigurationTemplateResource = metav1.APIResource{
		Name:         "applicationconfigurationtemplates",
		SingularName: "applicationconfigurationtemplate",
		Namespaced:   false,
		Kind:         ApplicationConfigurationTemplateGroupVersionKind.Kind,
	}
)

func NewApplicationConfigurationTemplate(namespace, name string, obj ApplicationConfigurationTemplate) *ApplicationConfigurationTemplate {
	obj.APIVersion, obj.Kind = ApplicationConfigurationTemplateGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type ApplicationConfigurationTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApplicationConfigurationTemplate
}

type ApplicationConfigurationTemplateHandlerFunc func(key string, obj *ApplicationConfigurationTemplate) (runtime.Object, error)

type ApplicationConfigurationTemplateChangeHandlerFunc func(obj *ApplicationConfigurationTemplate) (runtime.Object, error)

type ApplicationConfigurationTemplateLister interface {
	List(namespace string, selector labels.Selector) (ret []*ApplicationConfigurationTemplate, err error)
	Get(namespace, name string) (*ApplicationConfigurationTemplate, error)
}

type ApplicationConfigurationTemplateController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() ApplicationConfigurationTemplateLister
	AddHandler(ctx context.Context, name string, handler ApplicationConfigurationTemplateHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler ApplicationConfigurationTemplateHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type ApplicationConfigurationTemplateInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*ApplicationConfigurationTemplate) (*ApplicationConfigurationTemplate, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*ApplicationConfigurationTemplate, error)
	Get(name string, opts metav1.GetOptions) (*ApplicationConfigurationTemplate, error)
	Update(*ApplicationConfigurationTemplate) (*ApplicationConfigurationTemplate, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*ApplicationConfigurationTemplateList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() ApplicationConfigurationTemplateController
	AddHandler(ctx context.Context, name string, sync ApplicationConfigurationTemplateHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle ApplicationConfigurationTemplateLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ApplicationConfigurationTemplateHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ApplicationConfigurationTemplateLifecycle)
}

type applicationConfigurationTemplateLister struct {
	controller *applicationConfigurationTemplateController
}

func (l *applicationConfigurationTemplateLister) List(namespace string, selector labels.Selector) (ret []*ApplicationConfigurationTemplate, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*ApplicationConfigurationTemplate))
	})
	return
}

func (l *applicationConfigurationTemplateLister) Get(namespace, name string) (*ApplicationConfigurationTemplate, error) {
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
			Group:    ApplicationConfigurationTemplateGroupVersionKind.Group,
			Resource: "applicationConfigurationTemplate",
		}, key)
	}
	return obj.(*ApplicationConfigurationTemplate), nil
}

type applicationConfigurationTemplateController struct {
	controller.GenericController
}

func (c *applicationConfigurationTemplateController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *applicationConfigurationTemplateController) Lister() ApplicationConfigurationTemplateLister {
	return &applicationConfigurationTemplateLister{
		controller: c,
	}
}

func (c *applicationConfigurationTemplateController) AddHandler(ctx context.Context, name string, handler ApplicationConfigurationTemplateHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*ApplicationConfigurationTemplate); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *applicationConfigurationTemplateController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler ApplicationConfigurationTemplateHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*ApplicationConfigurationTemplate); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type applicationConfigurationTemplateFactory struct {
}

func (c applicationConfigurationTemplateFactory) Object() runtime.Object {
	return &ApplicationConfigurationTemplate{}
}

func (c applicationConfigurationTemplateFactory) List() runtime.Object {
	return &ApplicationConfigurationTemplateList{}
}

func (s *applicationConfigurationTemplateClient) Controller() ApplicationConfigurationTemplateController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.applicationConfigurationTemplateControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(ApplicationConfigurationTemplateGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &applicationConfigurationTemplateController{
		GenericController: genericController,
	}

	s.client.applicationConfigurationTemplateControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type applicationConfigurationTemplateClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   ApplicationConfigurationTemplateController
}

func (s *applicationConfigurationTemplateClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *applicationConfigurationTemplateClient) Create(o *ApplicationConfigurationTemplate) (*ApplicationConfigurationTemplate, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*ApplicationConfigurationTemplate), err
}

func (s *applicationConfigurationTemplateClient) Get(name string, opts metav1.GetOptions) (*ApplicationConfigurationTemplate, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*ApplicationConfigurationTemplate), err
}

func (s *applicationConfigurationTemplateClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*ApplicationConfigurationTemplate, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*ApplicationConfigurationTemplate), err
}

func (s *applicationConfigurationTemplateClient) Update(o *ApplicationConfigurationTemplate) (*ApplicationConfigurationTemplate, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*ApplicationConfigurationTemplate), err
}

func (s *applicationConfigurationTemplateClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *applicationConfigurationTemplateClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *applicationConfigurationTemplateClient) List(opts metav1.ListOptions) (*ApplicationConfigurationTemplateList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*ApplicationConfigurationTemplateList), err
}

func (s *applicationConfigurationTemplateClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *applicationConfigurationTemplateClient) Patch(o *ApplicationConfigurationTemplate, patchType types.PatchType, data []byte, subresources ...string) (*ApplicationConfigurationTemplate, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*ApplicationConfigurationTemplate), err
}

func (s *applicationConfigurationTemplateClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *applicationConfigurationTemplateClient) AddHandler(ctx context.Context, name string, sync ApplicationConfigurationTemplateHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *applicationConfigurationTemplateClient) AddLifecycle(ctx context.Context, name string, lifecycle ApplicationConfigurationTemplateLifecycle) {
	sync := NewApplicationConfigurationTemplateLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *applicationConfigurationTemplateClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ApplicationConfigurationTemplateHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *applicationConfigurationTemplateClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ApplicationConfigurationTemplateLifecycle) {
	sync := NewApplicationConfigurationTemplateLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type ApplicationConfigurationTemplateIndexer func(obj *ApplicationConfigurationTemplate) ([]string, error)

type ApplicationConfigurationTemplateClientCache interface {
	Get(namespace, name string) (*ApplicationConfigurationTemplate, error)
	List(namespace string, selector labels.Selector) ([]*ApplicationConfigurationTemplate, error)

	Index(name string, indexer ApplicationConfigurationTemplateIndexer)
	GetIndexed(name, key string) ([]*ApplicationConfigurationTemplate, error)
}

type ApplicationConfigurationTemplateClient interface {
	Create(*ApplicationConfigurationTemplate) (*ApplicationConfigurationTemplate, error)
	Get(namespace, name string, opts metav1.GetOptions) (*ApplicationConfigurationTemplate, error)
	Update(*ApplicationConfigurationTemplate) (*ApplicationConfigurationTemplate, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*ApplicationConfigurationTemplateList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() ApplicationConfigurationTemplateClientCache

	OnCreate(ctx context.Context, name string, sync ApplicationConfigurationTemplateChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync ApplicationConfigurationTemplateChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync ApplicationConfigurationTemplateChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() ApplicationConfigurationTemplateInterface
}

type applicationConfigurationTemplateClientCache struct {
	client *applicationConfigurationTemplateClient2
}

type applicationConfigurationTemplateClient2 struct {
	iface      ApplicationConfigurationTemplateInterface
	controller ApplicationConfigurationTemplateController
}

func (n *applicationConfigurationTemplateClient2) Interface() ApplicationConfigurationTemplateInterface {
	return n.iface
}

func (n *applicationConfigurationTemplateClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *applicationConfigurationTemplateClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *applicationConfigurationTemplateClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *applicationConfigurationTemplateClient2) Create(obj *ApplicationConfigurationTemplate) (*ApplicationConfigurationTemplate, error) {
	return n.iface.Create(obj)
}

func (n *applicationConfigurationTemplateClient2) Get(namespace, name string, opts metav1.GetOptions) (*ApplicationConfigurationTemplate, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *applicationConfigurationTemplateClient2) Update(obj *ApplicationConfigurationTemplate) (*ApplicationConfigurationTemplate, error) {
	return n.iface.Update(obj)
}

func (n *applicationConfigurationTemplateClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *applicationConfigurationTemplateClient2) List(namespace string, opts metav1.ListOptions) (*ApplicationConfigurationTemplateList, error) {
	return n.iface.List(opts)
}

func (n *applicationConfigurationTemplateClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *applicationConfigurationTemplateClientCache) Get(namespace, name string) (*ApplicationConfigurationTemplate, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *applicationConfigurationTemplateClientCache) List(namespace string, selector labels.Selector) ([]*ApplicationConfigurationTemplate, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *applicationConfigurationTemplateClient2) Cache() ApplicationConfigurationTemplateClientCache {
	n.loadController()
	return &applicationConfigurationTemplateClientCache{
		client: n,
	}
}

func (n *applicationConfigurationTemplateClient2) OnCreate(ctx context.Context, name string, sync ApplicationConfigurationTemplateChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &applicationConfigurationTemplateLifecycleDelegate{create: sync})
}

func (n *applicationConfigurationTemplateClient2) OnChange(ctx context.Context, name string, sync ApplicationConfigurationTemplateChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &applicationConfigurationTemplateLifecycleDelegate{update: sync})
}

func (n *applicationConfigurationTemplateClient2) OnRemove(ctx context.Context, name string, sync ApplicationConfigurationTemplateChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &applicationConfigurationTemplateLifecycleDelegate{remove: sync})
}

func (n *applicationConfigurationTemplateClientCache) Index(name string, indexer ApplicationConfigurationTemplateIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*ApplicationConfigurationTemplate); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *applicationConfigurationTemplateClientCache) GetIndexed(name, key string) ([]*ApplicationConfigurationTemplate, error) {
	var result []*ApplicationConfigurationTemplate
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*ApplicationConfigurationTemplate); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *applicationConfigurationTemplateClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type applicationConfigurationTemplateLifecycleDelegate struct {
	create ApplicationConfigurationTemplateChangeHandlerFunc
	update ApplicationConfigurationTemplateChangeHandlerFunc
	remove ApplicationConfigurationTemplateChangeHandlerFunc
}

func (n *applicationConfigurationTemplateLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *applicationConfigurationTemplateLifecycleDelegate) Create(obj *ApplicationConfigurationTemplate) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *applicationConfigurationTemplateLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *applicationConfigurationTemplateLifecycleDelegate) Remove(obj *ApplicationConfigurationTemplate) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *applicationConfigurationTemplateLifecycleDelegate) Updated(obj *ApplicationConfigurationTemplate) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
