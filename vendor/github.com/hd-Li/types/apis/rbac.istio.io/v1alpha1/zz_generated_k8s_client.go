package v1alpha1

import (
	"context"
	"sync"

	"github.com/rancher/norman/controller"
	"github.com/rancher/norman/objectclient"
	"github.com/rancher/norman/objectclient/dynamic"
	"github.com/rancher/norman/restwatch"
	"k8s.io/client-go/rest"
)

type (
	contextKeyType        struct{}
	contextClientsKeyType struct{}
)

type Interface interface {
	RESTClient() rest.Interface
	controller.Starter

	ClusterRbacConfigsGetter
	ServiceRolesGetter
	ServiceRoleBindingsGetter
}

type Clients struct {
	Interface Interface

	ClusterRbacConfig  ClusterRbacConfigClient
	ServiceRole        ServiceRoleClient
	ServiceRoleBinding ServiceRoleBindingClient
}

type Client struct {
	sync.Mutex
	restClient rest.Interface
	starters   []controller.Starter

	clusterRbacConfigControllers  map[string]ClusterRbacConfigController
	serviceRoleControllers        map[string]ServiceRoleController
	serviceRoleBindingControllers map[string]ServiceRoleBindingController
}

func Factory(ctx context.Context, config rest.Config) (context.Context, controller.Starter, error) {
	c, err := NewForConfig(config)
	if err != nil {
		return ctx, nil, err
	}

	cs := NewClientsFromInterface(c)

	ctx = context.WithValue(ctx, contextKeyType{}, c)
	ctx = context.WithValue(ctx, contextClientsKeyType{}, cs)
	return ctx, c, nil
}

func ClientsFrom(ctx context.Context) *Clients {
	return ctx.Value(contextClientsKeyType{}).(*Clients)
}

func From(ctx context.Context) Interface {
	return ctx.Value(contextKeyType{}).(Interface)
}

func NewClients(config rest.Config) (*Clients, error) {
	iface, err := NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return NewClientsFromInterface(iface), nil
}

func NewClientsFromInterface(iface Interface) *Clients {
	return &Clients{
		Interface: iface,

		ClusterRbacConfig: &clusterRbacConfigClient2{
			iface: iface.ClusterRbacConfigs(""),
		},
		ServiceRole: &serviceRoleClient2{
			iface: iface.ServiceRoles(""),
		},
		ServiceRoleBinding: &serviceRoleBindingClient2{
			iface: iface.ServiceRoleBindings(""),
		},
	}
}

func NewForConfig(config rest.Config) (Interface, error) {
	if config.NegotiatedSerializer == nil {
		config.NegotiatedSerializer = dynamic.NegotiatedSerializer
	}

	restClient, err := restwatch.UnversionedRESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &Client{
		restClient: restClient,

		clusterRbacConfigControllers:  map[string]ClusterRbacConfigController{},
		serviceRoleControllers:        map[string]ServiceRoleController{},
		serviceRoleBindingControllers: map[string]ServiceRoleBindingController{},
	}, nil
}

func (c *Client) RESTClient() rest.Interface {
	return c.restClient
}

func (c *Client) Sync(ctx context.Context) error {
	return controller.Sync(ctx, c.starters...)
}

func (c *Client) Start(ctx context.Context, threadiness int) error {
	return controller.Start(ctx, threadiness, c.starters...)
}

type ClusterRbacConfigsGetter interface {
	ClusterRbacConfigs(namespace string) ClusterRbacConfigInterface
}

func (c *Client) ClusterRbacConfigs(namespace string) ClusterRbacConfigInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &ClusterRbacConfigResource, ClusterRbacConfigGroupVersionKind, clusterRbacConfigFactory{})
	return &clusterRbacConfigClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type ServiceRolesGetter interface {
	ServiceRoles(namespace string) ServiceRoleInterface
}

func (c *Client) ServiceRoles(namespace string) ServiceRoleInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &ServiceRoleResource, ServiceRoleGroupVersionKind, serviceRoleFactory{})
	return &serviceRoleClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type ServiceRoleBindingsGetter interface {
	ServiceRoleBindings(namespace string) ServiceRoleBindingInterface
}

func (c *Client) ServiceRoleBindings(namespace string) ServiceRoleBindingInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &ServiceRoleBindingResource, ServiceRoleBindingGroupVersionKind, serviceRoleBindingFactory{})
	return &serviceRoleBindingClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}
