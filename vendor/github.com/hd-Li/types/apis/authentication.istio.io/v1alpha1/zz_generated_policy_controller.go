package v1alpha1

import (
	"context"

	"github.com/hd-Li/types/pkg/istio/apis/authentication/v1alpha1"
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
	PolicyGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "Policy",
	}
	PolicyResource = metav1.APIResource{
		Name:         "policies",
		SingularName: "policy",
		Namespaced:   true,

		Kind: PolicyGroupVersionKind.Kind,
	}
)

func NewPolicy(namespace, name string, obj v1alpha1.Policy) *v1alpha1.Policy {
	obj.APIVersion, obj.Kind = PolicyGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type PolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []v1alpha1.Policy
}

type PolicyHandlerFunc func(key string, obj *v1alpha1.Policy) (runtime.Object, error)

type PolicyChangeHandlerFunc func(obj *v1alpha1.Policy) (runtime.Object, error)

type PolicyLister interface {
	List(namespace string, selector labels.Selector) (ret []*v1alpha1.Policy, err error)
	Get(namespace, name string) (*v1alpha1.Policy, error)
}

type PolicyController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() PolicyLister
	AddHandler(ctx context.Context, name string, handler PolicyHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler PolicyHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type PolicyInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*v1alpha1.Policy) (*v1alpha1.Policy, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha1.Policy, error)
	Get(name string, opts metav1.GetOptions) (*v1alpha1.Policy, error)
	Update(*v1alpha1.Policy) (*v1alpha1.Policy, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*PolicyList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() PolicyController
	AddHandler(ctx context.Context, name string, sync PolicyHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle PolicyLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync PolicyHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle PolicyLifecycle)
}

type policyLister struct {
	controller *policyController
}

func (l *policyLister) List(namespace string, selector labels.Selector) (ret []*v1alpha1.Policy, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*v1alpha1.Policy))
	})
	return
}

func (l *policyLister) Get(namespace, name string) (*v1alpha1.Policy, error) {
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
			Group:    PolicyGroupVersionKind.Group,
			Resource: "policy",
		}, key)
	}
	return obj.(*v1alpha1.Policy), nil
}

type policyController struct {
	controller.GenericController
}

func (c *policyController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *policyController) Lister() PolicyLister {
	return &policyLister{
		controller: c,
	}
}

func (c *policyController) AddHandler(ctx context.Context, name string, handler PolicyHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha1.Policy); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *policyController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler PolicyHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v1alpha1.Policy); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type policyFactory struct {
}

func (c policyFactory) Object() runtime.Object {
	return &v1alpha1.Policy{}
}

func (c policyFactory) List() runtime.Object {
	return &PolicyList{}
}

func (s *policyClient) Controller() PolicyController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.policyControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(PolicyGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &policyController{
		GenericController: genericController,
	}

	s.client.policyControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type policyClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   PolicyController
}

func (s *policyClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *policyClient) Create(o *v1alpha1.Policy) (*v1alpha1.Policy, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*v1alpha1.Policy), err
}

func (s *policyClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.Policy, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*v1alpha1.Policy), err
}

func (s *policyClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v1alpha1.Policy, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*v1alpha1.Policy), err
}

func (s *policyClient) Update(o *v1alpha1.Policy) (*v1alpha1.Policy, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*v1alpha1.Policy), err
}

func (s *policyClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *policyClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *policyClient) List(opts metav1.ListOptions) (*PolicyList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*PolicyList), err
}

func (s *policyClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *policyClient) Patch(o *v1alpha1.Policy, patchType types.PatchType, data []byte, subresources ...string) (*v1alpha1.Policy, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*v1alpha1.Policy), err
}

func (s *policyClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *policyClient) AddHandler(ctx context.Context, name string, sync PolicyHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *policyClient) AddLifecycle(ctx context.Context, name string, lifecycle PolicyLifecycle) {
	sync := NewPolicyLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *policyClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync PolicyHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *policyClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle PolicyLifecycle) {
	sync := NewPolicyLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type PolicyIndexer func(obj *v1alpha1.Policy) ([]string, error)

type PolicyClientCache interface {
	Get(namespace, name string) (*v1alpha1.Policy, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha1.Policy, error)

	Index(name string, indexer PolicyIndexer)
	GetIndexed(name, key string) ([]*v1alpha1.Policy, error)
}

type PolicyClient interface {
	Create(*v1alpha1.Policy) (*v1alpha1.Policy, error)
	Get(namespace, name string, opts metav1.GetOptions) (*v1alpha1.Policy, error)
	Update(*v1alpha1.Policy) (*v1alpha1.Policy, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*PolicyList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() PolicyClientCache

	OnCreate(ctx context.Context, name string, sync PolicyChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync PolicyChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync PolicyChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() PolicyInterface
}

type policyClientCache struct {
	client *policyClient2
}

type policyClient2 struct {
	iface      PolicyInterface
	controller PolicyController
}

func (n *policyClient2) Interface() PolicyInterface {
	return n.iface
}

func (n *policyClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *policyClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *policyClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *policyClient2) Create(obj *v1alpha1.Policy) (*v1alpha1.Policy, error) {
	return n.iface.Create(obj)
}

func (n *policyClient2) Get(namespace, name string, opts metav1.GetOptions) (*v1alpha1.Policy, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *policyClient2) Update(obj *v1alpha1.Policy) (*v1alpha1.Policy, error) {
	return n.iface.Update(obj)
}

func (n *policyClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *policyClient2) List(namespace string, opts metav1.ListOptions) (*PolicyList, error) {
	return n.iface.List(opts)
}

func (n *policyClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *policyClientCache) Get(namespace, name string) (*v1alpha1.Policy, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *policyClientCache) List(namespace string, selector labels.Selector) ([]*v1alpha1.Policy, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *policyClient2) Cache() PolicyClientCache {
	n.loadController()
	return &policyClientCache{
		client: n,
	}
}

func (n *policyClient2) OnCreate(ctx context.Context, name string, sync PolicyChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &policyLifecycleDelegate{create: sync})
}

func (n *policyClient2) OnChange(ctx context.Context, name string, sync PolicyChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &policyLifecycleDelegate{update: sync})
}

func (n *policyClient2) OnRemove(ctx context.Context, name string, sync PolicyChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &policyLifecycleDelegate{remove: sync})
}

func (n *policyClientCache) Index(name string, indexer PolicyIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*v1alpha1.Policy); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *policyClientCache) GetIndexed(name, key string) ([]*v1alpha1.Policy, error) {
	var result []*v1alpha1.Policy
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*v1alpha1.Policy); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *policyClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type policyLifecycleDelegate struct {
	create PolicyChangeHandlerFunc
	update PolicyChangeHandlerFunc
	remove PolicyChangeHandlerFunc
}

func (n *policyLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *policyLifecycleDelegate) Create(obj *v1alpha1.Policy) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *policyLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *policyLifecycleDelegate) Remove(obj *v1alpha1.Policy) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *policyLifecycleDelegate) Updated(obj *v1alpha1.Policy) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
