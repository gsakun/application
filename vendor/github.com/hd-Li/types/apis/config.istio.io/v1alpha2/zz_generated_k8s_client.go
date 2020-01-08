package v1alpha2

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

	HandlersGetter
	InstancesGetter
	RulesGetter
	QuotaSpecsGetter
	QuotaSpecBindingsGetter
}

type Clients struct {
	Interface Interface

	Handler          HandlerClient
	Instance         InstanceClient
	Rule             RuleClient
	QuotaSpec        QuotaSpecClient
	QuotaSpecBinding QuotaSpecBindingClient
}

type Client struct {
	sync.Mutex
	restClient rest.Interface
	starters   []controller.Starter

	handlerControllers          map[string]HandlerController
	instanceControllers         map[string]InstanceController
	ruleControllers             map[string]RuleController
	quotaSpecControllers        map[string]QuotaSpecController
	quotaSpecBindingControllers map[string]QuotaSpecBindingController
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

		Handler: &handlerClient2{
			iface: iface.Handlers(""),
		},
		Instance: &instanceClient2{
			iface: iface.Instances(""),
		},
		Rule: &ruleClient2{
			iface: iface.Rules(""),
		},
		QuotaSpec: &quotaSpecClient2{
			iface: iface.QuotaSpecs(""),
		},
		QuotaSpecBinding: &quotaSpecBindingClient2{
			iface: iface.QuotaSpecBindings(""),
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

		handlerControllers:          map[string]HandlerController{},
		instanceControllers:         map[string]InstanceController{},
		ruleControllers:             map[string]RuleController{},
		quotaSpecControllers:        map[string]QuotaSpecController{},
		quotaSpecBindingControllers: map[string]QuotaSpecBindingController{},
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

type HandlersGetter interface {
	Handlers(namespace string) HandlerInterface
}

func (c *Client) Handlers(namespace string) HandlerInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &HandlerResource, HandlerGroupVersionKind, handlerFactory{})
	return &handlerClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type InstancesGetter interface {
	Instances(namespace string) InstanceInterface
}

func (c *Client) Instances(namespace string) InstanceInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &InstanceResource, InstanceGroupVersionKind, instanceFactory{})
	return &instanceClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type RulesGetter interface {
	Rules(namespace string) RuleInterface
}

func (c *Client) Rules(namespace string) RuleInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &RuleResource, RuleGroupVersionKind, ruleFactory{})
	return &ruleClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type QuotaSpecsGetter interface {
	QuotaSpecs(namespace string) QuotaSpecInterface
}

func (c *Client) QuotaSpecs(namespace string) QuotaSpecInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &QuotaSpecResource, QuotaSpecGroupVersionKind, quotaSpecFactory{})
	return &quotaSpecClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type QuotaSpecBindingsGetter interface {
	QuotaSpecBindings(namespace string) QuotaSpecBindingInterface
}

func (c *Client) QuotaSpecBindings(namespace string) QuotaSpecBindingInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &QuotaSpecBindingResource, QuotaSpecBindingGroupVersionKind, quotaSpecBindingFactory{})
	return &quotaSpecBindingClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}
