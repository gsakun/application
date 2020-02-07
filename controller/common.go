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

func GetOwnerRef(app *v3.Application) metav1.OwnerReference {
	ownerRef := metav1.OwnerReference{
		Name:       app.Namespace + ":" + app.Name,
		APIVersion: app.APIVersion,
		UID:        app.UID,
		Kind:       app.Kind,
	}

	return ownerRef
}

func GetOwnerRefFromNamespace(ns *corev1.Namespace) metav1.OwnerReference {
	ownerRef := metav1.OwnerReference{
		Name:       ns.Name,
		APIVersion: ns.APIVersion,
		UID:        ns.UID,
		Kind:       ns.Kind,
	}

	return ownerRef
}

func NewGatewayObject(app *v3.Application, ns *corev1.Namespace) istiov1alpha3.Gateway {
	ownerRef := GetOwnerRefFromNamespace(ns)

	gateway := istiov1alpha3.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Namespace + "-" + "gateway",
		},
		Spec: istiov1alpha3.GatewaySpec{
			Selector: map[string]string{"app": "istio-ingressgateway"},
			Servers: []istiov1alpha3.Server{
				istiov1alpha3.Server{
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

func NewPolicyObject(app *v3.Application, ns *corev1.Namespace) istioauthnv1alphav1.Policy {
	ownerRef := GetOwnerRefFromNamespace(ns)

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
			OwnerReferences: []metav1.OwnerReference{ownerRef},
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

func NewClusterRbacConfig(app *v3.Application, ns *corev1.Namespace) istiorbacv1alpha1.ClusterRbacConfig {
	var labels map[string]string = make(map[string]string)
	var ann map[string]string = make(map[string]string)
	labels[app.Namespace] = "included"
	rbacConfig := istiorbacv1alpha1.ClusterRbacConfig{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "istio-system",
			Name:      "default",
			//Labels:      map[string]string{app.Namespace: "included"},
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

func GetObjectApplied(obj interface{}) string {
	b, _ := json.Marshal(obj)
	//var out bytes.Buffer
	//json.Indent(&out, b, "", "")
	//return out.String()
	return string(b)
}
