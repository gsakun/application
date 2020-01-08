package v1alpha2

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
type AttributeManifest struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the implementation of this definition.
	// +optional
	Spec v1beta1.AttributeManifest `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// AttributeManifestList is a collection of AttributeManifests.
type AttributeManifestList struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items       []AttributeManifest `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type HTTPAPISpec struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the implementation of this definition.
	// +optional
	Spec client.HTTPAPISpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// HTTPAPISpecList is a collection of HTTPAPISpecs.
type HTTPAPISpecList struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items       []HTTPAPISpec `json:"items" protobuf:"bytes,2,rep,name=items"`
}


// HTTPAPISpecBinding defines the binding between HTTPAPISpecs and one or more
// IstioService. For example, the following establishes a binding
// between the HTTPAPISpec `petstore` and service `foo` in namespace `bar`.
type HTTPAPISpecBinding struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the implementation of this definition.
	// +optional
	Spec client.HTTPAPISpecBinding `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// HTTPAPISpecBindingList is a collection of HTTPAPISpecBindings.
type HTTPAPISpecBindingList struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items       []HTTPAPISpecBinding `json:"items" protobuf:"bytes,2,rep,name=items"`
}


// Each adapter implementation defines its own `params` proto.
//
// In the following example we define a `metrics` handler for the `prometheus` adapter.
// The example is in the form of a Kubernetes resource:
//  The `metadata.name` is the name of the handler
//  The `kind` refers to the adapter name
//  The `spec` block represents adapter-specific configuration as well as the connection information
//
// ```yaml
// # Sample-1: No connection specified (for compiled in adapters)
// # Note: if connection information is not specified, the adapter configuration is directly inside
// # `spec` block. This is going to be DEPRECATED in favor of Sample-2
// apiVersion: "config.istio.io/v1alpha2"
// kind: handler
// metadata:
//   name: requestcount
//   namespace: istio-system
// spec:
//   compiledAdapter: prometheus
//   params:
//     metrics:
//     - name: request_count
//       instance_name: requestcount.metric.istio-system
//       kind: COUNTER
//       label_names:
//       - source_service
//       - source_version
//       - destination_service
//       - destination_version
// ---
// # Sample-2: With connection information (for out-of-process adapters)
// # Note: Unlike sample-1, the adapter configuration is parallel to `connection` and is nested inside `param` block.
// apiVersion: "config.istio.io/v1alpha2"
// kind: handler
// metadata:
//   name: requestcount
//   namespace: istio-system
// spec:
//   compiledAdapter: prometheus
//   params:
//     param:
//       metrics:
//       - name: request_count
//         instance_name: requestcount.metric.istio-system
//         kind: COUNTER
//         label_names:
//         - source_service
//         - source_version
//         - destination_service
//         - destination_version
//     connection:
//       address: localhost:8090
// ---
// ```
*/

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Handler allows the operator to configure a specific adapter implementation.
type Handler struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the implementation of this definition.
	Spec HandlerSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type HandlerSpec struct {
	// Must be unique in the entire Mixer configuration. Used by [Actions][istio.policy.v1beta1.Action.handler]
	// to refer to this handler.
	Name string `json:"name,omitempty"`
	// The name of the compiled in adapter this handler instantiates. For referencing non compiled-in
	// adapters, use the `adapter` field instead.
	//
	// The value must match the name of the available adapter Mixer is built with. An adapter's name is typically a
	// constant in its code.
	CompiledAdapter string `json:"compiledAdapter,omitempty"`
	// The name of a specific adapter implementation. For referencing compiled-in
	// adapters, use the `compiled_adapter` field instead.
	//
	// An adapter's implementation name is typically a constant in its code.
	Adapter string `json:"adapter,omitempty"`
	// Depends on adapter implementation. Struct representation of a
	// proto defined by the adapter implementation; this varies depending on the value of field `adapter`.
	Params HandlerParams `json:"params,omitempty"`
	// Information on how to connect to the out-of-process adapter.
	// This is used if the adapter is not compiled into Mixer binary and is running as a separate process.
	//Connection *Connection `protobuf:"bytes,4,opt,name=connection,proto3" json:"connection,omitempty"`
}

type HandlerParams struct {
	RedisServerUrl string `json:"redisServerUrl,omitempty"`
	ConnectionPoolSize int32 `json:"connectionPoolSize,omitempty"`
	Quotas  []HandlerQuota  `json:"quotas,omitempty"`
}

type HandlerQuota struct {
	Name   string  `json:"name"`
	MaxAmount  int32  `json:"maxAmount"`
	ValidDuration string `json:"validDuration"`
	BucketDuration  string  `json:"bucketDuration"`
	RateLimitAlgorithm RateLimitAlgorithm `json:"rateLimitAlgorithm"`
	Overrides []Override `json:"overrides,omitempty"`
}

type Override struct{
	Dimensions map[string]string  `json:"dimensions"`
	MaxAmount  int32  `json:"maxAmount"`
}
type RateLimitAlgorithm string

const (
	ROLLING RateLimitAlgorithm = "ROLLING_WINDOW"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HandlerList is a collection of Handlers.
type HandlerList struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items       []Handler `json:"items" protobuf:"bytes,2,rep,name=items"`
}


// Instance is defined by the operator. Instance is defined relative to a known
// template. Their purpose is to tell Mixer how to use attributes or literals to produce
// instances of the specified template at runtime.
//
// The following example instructs Mixer to construct an instance associated with template
// 'istio.mixer.adapter.metric.Metric'. It provides a mapping from the template's fields to expressions.
// Instances produced with this instance can be referenced by [Actions][istio.policy.v1beta1.Action] using name
// 'RequestCountByService'
//
// ```yaml
// - name: RequestCountByService
//   template: istio.mixer.adapter.metric.Metric
//   params:
//     value: 1
//     dimensions:
//       source: source.name
//       destination_ip: destination.ip
// ```

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// An Instance tells Mixer how to create instances for particular template.
type Instance struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the implementation of this definition.
	Spec InstanceSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type InstanceSpec struct {
	// The name of this instance
	//
	// Must be unique amongst other Instances in scope. Used by [Action][istio.policy.v1beta1.Action] to refer
	// to an instance produced by this instance.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The name of the compiled in template this instance creates instances for. For referencing non compiled-in
	// templates, use the `template` field instead.
	//
	// The value must match the name of the available template Mixer is built with.
	CompiledTemplate string `protobuf:"bytes,67794676,opt,name=compiled_template,json=compiledTemplate,proto3" json:"compiled_template,omitempty"`
	// The name of the template this instance creates instances for. For referencing compiled-in
	// templates, use the `compiled_template` field instead.
	//
	// The value must match the name of the available template in scope.
	Template string `protobuf:"bytes,2,opt,name=template,proto3" json:"template,omitempty"`
	// Depends on referenced template. Struct representation of a
	// proto defined by the template; this varies depending on the value of field `template`.
	Params InstanceParams `protobuf:"bytes,3,opt,name=params,proto3" json:"params,omitempty"`
	// Defines attribute bindings to map the output of attribute-producing adapters back into
	// the attribute space. The variable `output` refers to the output template instance produced
	// by the adapter.
	// The following example derives `source.namespace` from `source.uid` in the context of Kubernetes:
	// ```yaml
	// params:
	//   # Pass the required attribute data to the adapter
	//   source_uid: source.uid | ""
	// attribute_bindings:
	//   # Fill the new attributes from the adapter produced output
	//   source.namespace: output.source_namespace
	// ```
	AttributeBindings map[string]string `protobuf:"bytes,4,rep,name=attribute_bindings,json=attributeBindings,proto3" json:"attribute_bindings,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type InstanceParams struct {
	Dimensions map[string]string  `json:"dimensions"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstanceList is a collection of Instances.
type InstanceList struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items       []Instance `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Determines the quotas used for individual requests.
type QuotaSpec struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the implementation of this definition.
	// +optional
	Spec QuotaSubSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type QuotaSubSpec struct {
	Rules []*QuotaRule `protobuf:"bytes,1,rep,name=rules,proto3" json:"rules,omitempty"`
}

type QuotaRule struct {
	// If empty, match all request.
	// If any of match is true, it is matched.
	Match []*AttributeMatch `protobuf:"bytes,1,rep,name=match,proto3" json:"match,omitempty"`
	// The list of quotas to charge.
	Quotas []*Quota `protobuf:"bytes,2,rep,name=quotas,proto3" json:"quotas,omitempty"`
}

type AttributeMatch struct {
	// Map of attribute names to StringMatch type.
	// Each map element specifies one condition to match.
	//
	// Example:
	//
	//   clause:
	//     source.uid:
	//       exact: SOURCE_UID
	//     request.http_method:
	//       exact: POST
	Clause map[string]*StringMatch `protobuf:"bytes,1,rep,name=clause,proto3" json:"clause,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type StringMatch struct {
	Exact    string `json:"exact,omitempty"`
	Prefix   string `json:"prefix,omitempty"`
	Regex    string `json:"regex,omitempty"`
}

type Quota struct {
	// The quota name to charge
	Quota string `protobuf:"bytes,1,opt,name=quota,proto3" json:"quota,omitempty"`
	// The quota amount to charge
	Charge int32 `protobuf:"varint,2,opt,name=charge,proto3" json:"charge,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QuotaSpecList is a collection of QuotaSpecs.
type QuotaSpecList struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items       []QuotaSpec `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type QuotaSpecBinding struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the implementation of this definition.
	// +optional
	Spec QuotaSpecBindingSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type QuotaSpecBindingSpec struct {
	// One or more services to map the listed QuotaSpec onto.
	Services []*IstioService `protobuf:"bytes,1,rep,name=services,proto3" json:"services,omitempty"`
	// One or more QuotaSpec references that should be mapped to
	// the specified service(s). The aggregate collection of match
	// conditions defined in the QuotaSpecs should not overlap.
	QuotaSpecs []*QuotaSpecBindingQuotaSpecReference `protobuf:"bytes,2,rep,name=quota_specs,json=quotaSpecs,proto3" json:"quotaSpecs,omitempty"`
}

type IstioService struct {
	// The short name of the service such as "foo".
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Optional namespace of the service. Defaults to value of metadata namespace field.
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// Domain suffix used to construct the service FQDN in implementations that support such specification.
	Domain string `protobuf:"bytes,3,opt,name=domain,proto3" json:"domain,omitempty"`
	// The service FQDN.
	Service string `protobuf:"bytes,4,opt,name=service,proto3" json:"service,omitempty"`
	// Optional one or more labels that uniquely identify the service version.
	//
	// *Note:* When used for a VirtualService destination, labels MUST be empty.
	//
	Labels map[string]string `protobuf:"bytes,5,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type QuotaSpecBindingQuotaSpecReference struct {
	// The short name of the QuotaSpec. This is the resource
	// name defined by the metadata name field.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Optional namespace of the QuotaSpec. Defaults to the value of the
	// metadata namespace field.
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QuotaSpecBindingList is a collection of QuotaSpecBindings.
type QuotaSpecBindingList struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items       []QuotaSpecBinding `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// The following example instructs Mixer to invoke `prometheus-handler` handler for all services and pass it the
// instance constructed using the 'RequestCountByService' instance.
//
// ```yaml
// - match: match(destination.service.host, "*")
//   actions:
//   - handler: prometheus-handler
//     instances:
//     - RequestCountByService
// ```

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// A Rule is a selector and a set of intentions to be executed when the selector is `true`
type Rule struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the implementation of this definition.
	Spec RuleSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type RuleSpec struct {
	Match string `protobuf:"bytes,1,opt,name=match,proto3" json:"match,omitempty"`
	Actions []*Action `protobuf:"bytes,2,rep,name=actions,proto3" json:"actions,omitempty"`
}

type Action struct {
	// Fully qualified name of the handler to invoke.
	// Must match the `name` of a [Handler][istio.policy.v1beta1.Handler.name].
	Handler string `protobuf:"bytes,2,opt,name=handler,proto3" json:"handler,omitempty"`
	// Each value must match the fully qualified name of the
	// [Instance][istio.policy.v1beta1.Instance.name]s.
	// Referenced instances are evaluated by resolving the attributes/literals for all the fields.
	// The constructed objects are then passed to the `handler` referenced within this action.
	Instances []string `protobuf:"bytes,3,rep,name=instances,proto3" json:"instances,omitempty"`
	// A handle to refer to the results of the action.
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RuleList is a collection of Rules.
type RuleList struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items       []Rule `json:"items" protobuf:"bytes,2,rep,name=items"`
}
