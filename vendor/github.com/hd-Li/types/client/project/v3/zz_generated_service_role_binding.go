package client

import (
	"github.com/rancher/norman/types"
)

const (
	ServiceRoleBindingType                      = "serviceRoleBinding"
	ServiceRoleBindingFieldAnnotations          = "annotations"
	ServiceRoleBindingFieldCreated              = "created"
	ServiceRoleBindingFieldCreatorID            = "creatorId"
	ServiceRoleBindingFieldLabels               = "labels"
	ServiceRoleBindingFieldMode                 = "mode"
	ServiceRoleBindingFieldName                 = "name"
	ServiceRoleBindingFieldNamespaceId          = "namespaceId"
	ServiceRoleBindingFieldOwnerReferences      = "ownerReferences"
	ServiceRoleBindingFieldProjectID            = "projectId"
	ServiceRoleBindingFieldRemoved              = "removed"
	ServiceRoleBindingFieldRole                 = "role"
	ServiceRoleBindingFieldRoleRef              = "roleRef"
	ServiceRoleBindingFieldState                = "state"
	ServiceRoleBindingFieldStatus               = "status"
	ServiceRoleBindingFieldSubjects             = "subjects"
	ServiceRoleBindingFieldTransitioning        = "transitioning"
	ServiceRoleBindingFieldTransitioningMessage = "transitioningMessage"
	ServiceRoleBindingFieldUUID                 = "uuid"
)

type ServiceRoleBinding struct {
	types.Resource
	Annotations          map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Created              string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID            string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Labels               map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Mode                 int64             `json:"mode,omitempty" yaml:"mode,omitempty"`
	Name                 string            `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceId          string            `json:"namespaceId,omitempty" yaml:"namespaceId,omitempty"`
	OwnerReferences      []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	ProjectID            string            `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	Removed              string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	Role                 string            `json:"role,omitempty" yaml:"role,omitempty"`
	RoleRef              *RoleRef          `json:"roleRef,omitempty" yaml:"roleRef,omitempty"`
	State                string            `json:"state,omitempty" yaml:"state,omitempty"`
	Status               interface{}       `json:"status,omitempty" yaml:"status,omitempty"`
	Subjects             []Subject         `json:"subjects,omitempty" yaml:"subjects,omitempty"`
	Transitioning        string            `json:"transitioning,omitempty" yaml:"transitioning,omitempty"`
	TransitioningMessage string            `json:"transitioningMessage,omitempty" yaml:"transitioningMessage,omitempty"`
	UUID                 string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type ServiceRoleBindingCollection struct {
	types.Collection
	Data   []ServiceRoleBinding `json:"data,omitempty"`
	client *ServiceRoleBindingClient
}

type ServiceRoleBindingClient struct {
	apiClient *Client
}

type ServiceRoleBindingOperations interface {
	List(opts *types.ListOpts) (*ServiceRoleBindingCollection, error)
	Create(opts *ServiceRoleBinding) (*ServiceRoleBinding, error)
	Update(existing *ServiceRoleBinding, updates interface{}) (*ServiceRoleBinding, error)
	Replace(existing *ServiceRoleBinding) (*ServiceRoleBinding, error)
	ByID(id string) (*ServiceRoleBinding, error)
	Delete(container *ServiceRoleBinding) error
}

func newServiceRoleBindingClient(apiClient *Client) *ServiceRoleBindingClient {
	return &ServiceRoleBindingClient{
		apiClient: apiClient,
	}
}

func (c *ServiceRoleBindingClient) Create(container *ServiceRoleBinding) (*ServiceRoleBinding, error) {
	resp := &ServiceRoleBinding{}
	err := c.apiClient.Ops.DoCreate(ServiceRoleBindingType, container, resp)
	return resp, err
}

func (c *ServiceRoleBindingClient) Update(existing *ServiceRoleBinding, updates interface{}) (*ServiceRoleBinding, error) {
	resp := &ServiceRoleBinding{}
	err := c.apiClient.Ops.DoUpdate(ServiceRoleBindingType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *ServiceRoleBindingClient) Replace(obj *ServiceRoleBinding) (*ServiceRoleBinding, error) {
	resp := &ServiceRoleBinding{}
	err := c.apiClient.Ops.DoReplace(ServiceRoleBindingType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *ServiceRoleBindingClient) List(opts *types.ListOpts) (*ServiceRoleBindingCollection, error) {
	resp := &ServiceRoleBindingCollection{}
	err := c.apiClient.Ops.DoList(ServiceRoleBindingType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *ServiceRoleBindingCollection) Next() (*ServiceRoleBindingCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &ServiceRoleBindingCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *ServiceRoleBindingClient) ByID(id string) (*ServiceRoleBinding, error) {
	resp := &ServiceRoleBinding{}
	err := c.apiClient.Ops.DoByID(ServiceRoleBindingType, id, resp)
	return resp, err
}

func (c *ServiceRoleBindingClient) Delete(container *ServiceRoleBinding) error {
	return c.apiClient.Ops.DoResourceDelete(ServiceRoleBindingType, &container.Resource)
}
