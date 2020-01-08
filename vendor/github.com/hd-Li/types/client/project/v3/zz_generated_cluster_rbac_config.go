package client

import (
	"github.com/rancher/norman/types"
)

const (
	ClusterRbacConfigType                      = "clusterRbacConfig"
	ClusterRbacConfigFieldAnnotations          = "annotations"
	ClusterRbacConfigFieldCreated              = "created"
	ClusterRbacConfigFieldCreatorID            = "creatorId"
	ClusterRbacConfigFieldEnforcementMode      = "enforcementMode"
	ClusterRbacConfigFieldExclusion            = "exclusion"
	ClusterRbacConfigFieldInclusion            = "inclusion"
	ClusterRbacConfigFieldLabels               = "labels"
	ClusterRbacConfigFieldMode                 = "mode"
	ClusterRbacConfigFieldName                 = "name"
	ClusterRbacConfigFieldNamespaceId          = "namespaceId"
	ClusterRbacConfigFieldOwnerReferences      = "ownerReferences"
	ClusterRbacConfigFieldProjectID            = "projectId"
	ClusterRbacConfigFieldRemoved              = "removed"
	ClusterRbacConfigFieldState                = "state"
	ClusterRbacConfigFieldStatus               = "status"
	ClusterRbacConfigFieldTransitioning        = "transitioning"
	ClusterRbacConfigFieldTransitioningMessage = "transitioningMessage"
	ClusterRbacConfigFieldUUID                 = "uuid"
)

type ClusterRbacConfig struct {
	types.Resource
	Annotations          map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Created              string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID            string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	EnforcementMode      int64             `json:"enforcementMode,omitempty" yaml:"enforcementMode,omitempty"`
	Exclusion            *RbacConfigTarget `json:"exclusion,omitempty" yaml:"exclusion,omitempty"`
	Inclusion            *RbacConfigTarget `json:"inclusion,omitempty" yaml:"inclusion,omitempty"`
	Labels               map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Mode                 int64             `json:"mode,omitempty" yaml:"mode,omitempty"`
	Name                 string            `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceId          string            `json:"namespaceId,omitempty" yaml:"namespaceId,omitempty"`
	OwnerReferences      []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	ProjectID            string            `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	Removed              string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	State                string            `json:"state,omitempty" yaml:"state,omitempty"`
	Status               interface{}       `json:"status,omitempty" yaml:"status,omitempty"`
	Transitioning        string            `json:"transitioning,omitempty" yaml:"transitioning,omitempty"`
	TransitioningMessage string            `json:"transitioningMessage,omitempty" yaml:"transitioningMessage,omitempty"`
	UUID                 string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type ClusterRbacConfigCollection struct {
	types.Collection
	Data   []ClusterRbacConfig `json:"data,omitempty"`
	client *ClusterRbacConfigClient
}

type ClusterRbacConfigClient struct {
	apiClient *Client
}

type ClusterRbacConfigOperations interface {
	List(opts *types.ListOpts) (*ClusterRbacConfigCollection, error)
	Create(opts *ClusterRbacConfig) (*ClusterRbacConfig, error)
	Update(existing *ClusterRbacConfig, updates interface{}) (*ClusterRbacConfig, error)
	Replace(existing *ClusterRbacConfig) (*ClusterRbacConfig, error)
	ByID(id string) (*ClusterRbacConfig, error)
	Delete(container *ClusterRbacConfig) error
}

func newClusterRbacConfigClient(apiClient *Client) *ClusterRbacConfigClient {
	return &ClusterRbacConfigClient{
		apiClient: apiClient,
	}
}

func (c *ClusterRbacConfigClient) Create(container *ClusterRbacConfig) (*ClusterRbacConfig, error) {
	resp := &ClusterRbacConfig{}
	err := c.apiClient.Ops.DoCreate(ClusterRbacConfigType, container, resp)
	return resp, err
}

func (c *ClusterRbacConfigClient) Update(existing *ClusterRbacConfig, updates interface{}) (*ClusterRbacConfig, error) {
	resp := &ClusterRbacConfig{}
	err := c.apiClient.Ops.DoUpdate(ClusterRbacConfigType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *ClusterRbacConfigClient) Replace(obj *ClusterRbacConfig) (*ClusterRbacConfig, error) {
	resp := &ClusterRbacConfig{}
	err := c.apiClient.Ops.DoReplace(ClusterRbacConfigType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *ClusterRbacConfigClient) List(opts *types.ListOpts) (*ClusterRbacConfigCollection, error) {
	resp := &ClusterRbacConfigCollection{}
	err := c.apiClient.Ops.DoList(ClusterRbacConfigType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *ClusterRbacConfigCollection) Next() (*ClusterRbacConfigCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &ClusterRbacConfigCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *ClusterRbacConfigClient) ByID(id string) (*ClusterRbacConfig, error) {
	resp := &ClusterRbacConfig{}
	err := c.apiClient.Ops.DoByID(ClusterRbacConfigType, id, resp)
	return resp, err
}

func (c *ClusterRbacConfigClient) Delete(container *ClusterRbacConfig) error {
	return c.apiClient.Ops.DoResourceDelete(ClusterRbacConfigType, &container.Resource)
}
