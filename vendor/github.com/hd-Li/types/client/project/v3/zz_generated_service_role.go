package client

import (
	"github.com/rancher/norman/types"
)

const (
	ServiceRoleType                      = "serviceRole"
	ServiceRoleFieldAnnotations          = "annotations"
	ServiceRoleFieldCreated              = "created"
	ServiceRoleFieldCreatorID            = "creatorId"
	ServiceRoleFieldLabels               = "labels"
	ServiceRoleFieldName                 = "name"
	ServiceRoleFieldNamespaceId          = "namespaceId"
	ServiceRoleFieldOwnerReferences      = "ownerReferences"
	ServiceRoleFieldProjectID            = "projectId"
	ServiceRoleFieldRemoved              = "removed"
	ServiceRoleFieldRules                = "rules"
	ServiceRoleFieldState                = "state"
	ServiceRoleFieldStatus               = "status"
	ServiceRoleFieldTransitioning        = "transitioning"
	ServiceRoleFieldTransitioningMessage = "transitioningMessage"
	ServiceRoleFieldUUID                 = "uuid"
)

type ServiceRole struct {
	types.Resource
	Annotations          map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Created              string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID            string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Labels               map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name                 string            `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceId          string            `json:"namespaceId,omitempty" yaml:"namespaceId,omitempty"`
	OwnerReferences      []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	ProjectID            string            `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	Removed              string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	Rules                []AccessRule      `json:"rules,omitempty" yaml:"rules,omitempty"`
	State                string            `json:"state,omitempty" yaml:"state,omitempty"`
	Status               interface{}       `json:"status,omitempty" yaml:"status,omitempty"`
	Transitioning        string            `json:"transitioning,omitempty" yaml:"transitioning,omitempty"`
	TransitioningMessage string            `json:"transitioningMessage,omitempty" yaml:"transitioningMessage,omitempty"`
	UUID                 string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type ServiceRoleCollection struct {
	types.Collection
	Data   []ServiceRole `json:"data,omitempty"`
	client *ServiceRoleClient
}

type ServiceRoleClient struct {
	apiClient *Client
}

type ServiceRoleOperations interface {
	List(opts *types.ListOpts) (*ServiceRoleCollection, error)
	Create(opts *ServiceRole) (*ServiceRole, error)
	Update(existing *ServiceRole, updates interface{}) (*ServiceRole, error)
	Replace(existing *ServiceRole) (*ServiceRole, error)
	ByID(id string) (*ServiceRole, error)
	Delete(container *ServiceRole) error
}

func newServiceRoleClient(apiClient *Client) *ServiceRoleClient {
	return &ServiceRoleClient{
		apiClient: apiClient,
	}
}

func (c *ServiceRoleClient) Create(container *ServiceRole) (*ServiceRole, error) {
	resp := &ServiceRole{}
	err := c.apiClient.Ops.DoCreate(ServiceRoleType, container, resp)
	return resp, err
}

func (c *ServiceRoleClient) Update(existing *ServiceRole, updates interface{}) (*ServiceRole, error) {
	resp := &ServiceRole{}
	err := c.apiClient.Ops.DoUpdate(ServiceRoleType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *ServiceRoleClient) Replace(obj *ServiceRole) (*ServiceRole, error) {
	resp := &ServiceRole{}
	err := c.apiClient.Ops.DoReplace(ServiceRoleType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *ServiceRoleClient) List(opts *types.ListOpts) (*ServiceRoleCollection, error) {
	resp := &ServiceRoleCollection{}
	err := c.apiClient.Ops.DoList(ServiceRoleType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *ServiceRoleCollection) Next() (*ServiceRoleCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &ServiceRoleCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *ServiceRoleClient) ByID(id string) (*ServiceRole, error) {
	resp := &ServiceRole{}
	err := c.apiClient.Ops.DoByID(ServiceRoleType, id, resp)
	return resp, err
}

func (c *ServiceRoleClient) Delete(container *ServiceRole) error {
	return c.apiClient.Ops.DoResourceDelete(ServiceRoleType, &container.Resource)
}
