package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Disable Istio RBAC completely, Istio RBAC policies will not be enforced.
	OFF RbacConfigMode = 0
	// Enable Istio RBAC for all services and namespaces. Note Istio RBAC is deny-by-default
	// which means all requests will be denied if it's not allowed by RBAC rules.
	ON RbacConfigMode = 1
	// Enable Istio RBAC only for services and namespaces specified in the inclusion field. Any other
	// services and namespaces not in the inclusion field will not be enforced by Istio RBAC policies.
	ON_WITH_INCLUSION RbacConfigMode = 2
	// Enable Istio RBAC for all services and namespaces except those specified in the exclusion field. Any other
	// services and namespaces not in the exclusion field will be enforced by Istio RBAC policies.
	ON_WITH_EXCLUSION RbacConfigMode = 3
)

type RbacConfigMode int32

type ClusterRbacConfig struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec RbacConfigSpec `json:"spec,omitempty"`
}

type RbacConfigSpec struct {
	// Istio RBAC mode.
	Mode RbacConfigMode `json:"mode,omitempty"`
	// A list of services or namespaces that should be enforced by Istio RBAC policies. Note: This field have
	// effect only when mode is ON_WITH_INCLUSION and will be ignored for any other modes.
	Inclusion *RbacConfigTarget `json:"inclusion,omitempty"`
	// A list of services or namespaces that should not be enforced by Istio RBAC policies. Note: This field have
	// effect only when mode is ON_WITH_EXCLUSION and will be ignored for any other modes.
	Exclusion *RbacConfigTarget `json:"exclusion,omitempty"`

	EnforcementMode EnforcementMode `json:"enforcementMode,omitempty"`
}

// Target defines a list of services or namespaces.
type RbacConfigTarget struct {
	// A list of services.
	Services []string `json:"services,omitempty"`
	// A list of namespaces.
	Namespaces []string `json:"namespaces,omitempty"`
}

type EnforcementMode int32

const (
	// Policy in ENFORCED mode has impact on user experience.
	// Policy is in ENFORCED mode by default.
	ENFORCED EnforcementMode = 0
	// Policy in PERMISSIVE mode isn't enforced and has no impact on users.
	// RBAC engine run policies in PERMISSIVE mode and logs stats.
	PERMISSIVE EnforcementMode = 1
)


type RbacConfig struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the implementation of this definition.
	// +optional
	Spec RbacConfigSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// RbacConfigList is a collection of RbacConfigs.
type RbacConfigList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items       []RbacConfig `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// ClusterRbacConfigList is a collection of ClusterRbacConfigs.
type ClusterRbacConfigList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items       []ClusterRbacConfig `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// ServiceRole specification contains a list of access rules (permissions).

type ServiceRole struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the implementation of this definition.
	// +optional
	Spec ServiceRoleSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type ServiceRoleSpec struct {
	Rules []AccessRule `protobuf:"bytes,1,rep,name=rules,proto3" json:"rules,omitempty"`
}

type AccessRule struct {
	// A list of service names.
	// Exact match, prefix match, and suffix match are supported for service names.
	// For example, the service name "bookstore.mtv.cluster.local" matches
	// "bookstore.mtv.cluster.local" (exact match), or "bookstore\*" (prefix match),
	// or "\*.mtv.cluster.local" (suffix match).
	// If set to ["\*"], it refers to all services in the namespace.
	Services []string `protobuf:"bytes,1,rep,name=services,proto3" json:"services,omitempty"`
	// $hide_from_docs
	// Optional. A list of HTTP hosts. This is matched against the HOST header in
	// a HTTP request. Exact match, prefix match and suffix match are supported.
	// For example, the host "test.abc.com" matches "test.abc.com" (exact match),
	// or "\*.abc.com" (prefix match), or "test.abc.\*" (suffix match).
	// If not specified, it matches to any host.
	// This field should not be set for TCP services. The policy will be ignored.
	Hosts []string `protobuf:"bytes,5,rep,name=hosts,proto3" json:"hosts,omitempty"`
	// $hide_from_docs
	// Optional. A list of HTTP hosts that must not be matched.
	NotHosts []string `protobuf:"bytes,6,rep,name=not_hosts,json=notHosts,proto3" json:"not_hosts,omitempty"`
	// Optional. A list of HTTP paths or gRPC methods.
	// gRPC methods must be presented as fully-qualified name in the form of
	// "/packageName.serviceName/methodName" and are case sensitive.
	// Exact match, prefix match, and suffix match are supported. For example,
	// the path "/books/review" matches "/books/review" (exact match),
	// or "/books/\*" (prefix match), or "review" (suffix match).
	// If not specified, it matches to any path.
	// This field should not be set for TCP services. The policy will be ignored.
	Paths []string `protobuf:"bytes,2,rep,name=paths,proto3" json:"paths,omitempty"`
	// $hide_from_docs
	// Optional. A list of HTTP paths or gRPC methods that must not be matched.
	NotPaths []string `protobuf:"bytes,7,rep,name=not_paths,json=notPaths,proto3" json:"not_paths,omitempty"`
	// Optional. A list of HTTP methods (e.g., "GET", "POST").
	// If not specified or specified as "\*", it matches to any methods.
	// This field should not be set for TCP services. The policy will be ignored.
	// For gRPC services, only `POST` is allowed; other methods will result in denying services.
	Methods []string `protobuf:"bytes,3,rep,name=methods,proto3" json:"methods,omitempty"`
	// $hide_from_docs
	// Optional. A list of HTTP methods that must not be matched.
	// Note: It's an error to set methods and not_methods at the same time.
	NotMethods []string `protobuf:"bytes,8,rep,name=not_methods,json=notMethods,proto3" json:"not_methods,omitempty"`
	// $hide_from_docs
	// Optional. A list of port numbers of the request. If not specified, it matches
	// to any port number.
	// Note: It's an error to set ports and not_ports at the same time.
	Ports []int32 `protobuf:"varint,9,rep,packed,name=ports,proto3" json:"ports,omitempty"`
	// $hide_from_docs
	// Optional.  A list of port numbers that must not be matched.
	// Note: It's an error to set ports and not_ports at the same time.
	NotPorts []int32 `protobuf:"varint,10,rep,packed,name=not_ports,json=notPorts,proto3" json:"not_ports,omitempty"`
	// Optional. Extra constraints in the ServiceRole specification.
	Constraints []AccessRule_Constraint `protobuf:"bytes,4,rep,name=constraints,proto3" json:"constraints,omitempty"`
}

type AccessRule_Constraint struct {
	// Key of the constraint.
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// List of valid values for the constraint.
	// Exact match, prefix match, and suffix match are supported.
	// For example, the value "v1alpha2" matches "v1alpha2" (exact match),
	// or "v1\*" (prefix match), or "\*alpha2" (suffix match).
	Values []string `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
}

// ServiceRoleList is a collection of ServiceRoles.
type ServiceRoleList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items       []ServiceRole `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type ServiceRoleBinding struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the implementation of this definition.
	// +optional
	Spec ServiceRoleBindingSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// ServiceRoleBindingList is a collection of ServiceRoleBindings.
type ServiceRoleBindingList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items       []ServiceRoleBinding `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type ServiceRoleBindingSpec struct {
	// List of subjects that are assigned the ServiceRole object.
	Subjects []Subject `protobuf:"bytes,1,rep,name=subjects,proto3" json:"subjects,omitempty"`
	// Reference to the ServiceRole object.
	RoleRef RoleRef `protobuf:"bytes,2,opt,name=roleRef,proto3" json:"roleRef,omitempty"`
	// $hide_from_docs
	// Indicates enforcement mode of the ServiceRoleBinding.
	Mode EnforcementMode `protobuf:"varint,3,opt,name=mode,proto3,enum=istio.rbac.v1alpha1.EnforcementMode" json:"mode,omitempty"`

	Actions []AccessRule `protobuf:"bytes,4,rep,name=actions,proto3" json:"actions,omitempty"`

	Role string `protobuf:"bytes,5,opt,name=role,proto3" json:"role,omitempty"`
}

type Subject struct {
	// Optional. The user name/ID that the subject represents.
	User string `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	// $hide_from_docs
	// Optional. A list of subject names. This is matched to the
	// `source.principal` attribute. If one of subject names is "\*", it matches to a subject with any name.
	// Prefix and suffix matches are supported.
	Names []string `protobuf:"bytes,4,rep,name=names,proto3" json:"names,omitempty"`
	// $hide_from_docs
	// Optional. A list of subject names that must not be matched.
	NotNames []string `protobuf:"bytes,5,rep,name=not_names,json=notNames,proto3" json:"not_names,omitempty"`
	// $hide_from_docs
	// Optional. The group that the subject belongs to.
	// Deprecated. Use groups and not_groups instead.
	Group string `protobuf:"bytes,2,opt,name=group,proto3" json:"group,omitempty"` // Deprecated: Do not use.
	// $hide_from_docs
	// Optional. A list of groups that the subject represents. This is matched to the
	// `request.auth.claims[groups]` attribute. If not specified, it applies to any groups.
	Groups []string `protobuf:"bytes,6,rep,name=groups,proto3" json:"groups,omitempty"`
	// $hide_from_docs
	// Optional. A list of groups that must not be matched.
	NotGroups []string `protobuf:"bytes,7,rep,name=not_groups,json=notGroups,proto3" json:"not_groups,omitempty"`
	// $hide_from_docs
	// Optional. A list of namespaces that the subject represents. This is matched to
	// the `source.namespace` attribute. If not specified, it applies to any namespaces.
	Namespaces []string `protobuf:"bytes,8,rep,name=namespaces,proto3" json:"namespaces,omitempty"`
	// $hide_from_docs
	// Optional. A list of namespaces that must not be matched.
	NotNamespaces []string `protobuf:"bytes,9,rep,name=not_namespaces,json=notNamespaces,proto3" json:"not_namespaces,omitempty"`
	// $hide_from_docs
	// Optional. A list of IP address or CIDR ranges that the subject represents.
	// E.g. 192.168.100.2 or 10.1.0.0/16. If not specified, it applies to any IP addresses.
	Ips []string `protobuf:"bytes,10,rep,name=ips,proto3" json:"ips,omitempty"`
	// $hide_from_docs
	// Optional. A list of IP addresses or CIDR ranges that must not be matched.
	NotIps []string `protobuf:"bytes,11,rep,name=not_ips,json=notIps,proto3" json:"not_ips,omitempty"`
	// Optional. The set of properties that identify the subject.
	Properties map[string]string `protobuf:"bytes,3,rep,name=properties,proto3" json:"properties,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type RoleRef struct {
	// The type of the role being referenced.
	// Currently, "ServiceRole" is the only supported value for "kind".
	Kind string `protobuf:"bytes,1,opt,name=kind,proto3" json:"kind,omitempty"`
	// The name of the ServiceRole object being referenced.
	// The ServiceRole object must be in the same namespace as the ServiceRoleBinding object.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}
