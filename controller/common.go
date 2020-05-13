package controller

import (
	"os"
	//"bytes"
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	//"k8s.io/apimachinery/pkg/runtime"

	v3 "github.com/hd-Li/types/apis/project.cattle.io/v3"
	istioauthnv1alphav1 "github.com/hd-Li/types/pkg/istio/apis/authentication/v1alpha1"
	istiorbacv1alpha1 "github.com/hd-Li/types/pkg/istio/apis/rbac/v1alpha1"
	istiov1alpha3 "github.com/knative/pkg/apis/istio/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewGatewayObject use for generate GatewayObject
func NewGatewayObject(app *v3.Application, ns *corev1.Namespace) istiov1alpha3.Gateway {
	gateway := istiov1alpha3.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(ns, corev1.SchemeGroupVersion.WithKind("Namespace"))},
			Namespace:       app.Namespace,
			Name:            app.Namespace + "-" + "gateway",
		},
		Spec: istiov1alpha3.GatewaySpec{
			Selector: map[string]string{"app": "istio-ingressgateway"},
			Servers: []istiov1alpha3.Server{
				{
					Hosts: []string{"*"},
					Port: istiov1alpha3.Port{
						Name:     "http",
						Number:   80,
						Protocol: istiov1alpha3.ProtocolHTTP,
					},
				},
			},
		},
	}

	return gateway
}

// NewPolicyObject use for generate PolicyObject
func NewPolicyObject(app *v3.Application, ns *corev1.Namespace) istioauthnv1alphav1.Policy {
	authnEndpoint := os.Getenv("AUTHN_ENDPOINT")
	realm := os.Getenv("AUTHN_REALM")

	issuer := authnEndpoint + "/auth/realms/" + realm
	uri := issuer + "/protocol/openid-connect/certs"

	originAuthenticationMethod := istioauthnv1alphav1.OriginAuthenticationMethod{
		Jwt: &istioauthnv1alphav1.Jwt{
			Issuer:  issuer,
			JwksUri: uri,
		},
	}

	policy := istioauthnv1alphav1.Policy{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(ns, corev1.SchemeGroupVersion.WithKind("Namespace"))},
			Namespace:       app.Namespace,
			Name:            "default",
		},
		Spec: istioauthnv1alphav1.PolicySpec{
			Origins:          []istioauthnv1alphav1.OriginAuthenticationMethod{originAuthenticationMethod},
			PrincipalBinding: istioauthnv1alphav1.USE_ORIGIN,
		},
	}

	return policy
}

// NewClusterRbacConfig use for generate ClusterRbacConfig
func NewClusterRbacConfig(app *v3.Application, ns *corev1.Namespace) istiorbacv1alpha1.ClusterRbacConfig {
	var labels map[string]string = make(map[string]string)
	var ann map[string]string = make(map[string]string)
	labels[app.Namespace] = "included"
	rbacConfig := istiorbacv1alpha1.ClusterRbacConfig{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   "default",
			Name:        "default",
			Labels:      labels,
			Annotations: ann,
		},
		Spec: istiorbacv1alpha1.RbacConfigSpec{
			Mode: istiorbacv1alpha1.ON_WITH_INCLUSION,
			Inclusion: &istiorbacv1alpha1.RbacConfigTarget{
				Namespaces: []string{app.Namespace},
			},
		},
	}

	return rbacConfig
}

// GetObjectApplied use for generate ObjectApplied
func GetObjectApplied(obj interface{}) string {
	b, _ := json.Marshal(obj)
	return string(b)
}
