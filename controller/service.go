package controller

import (
	"github.com/hd-Li/types/apis/project.cattle.io/v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	istiov1alpha3 "github.com/knative/pkg/apis/istio/v1alpha3"
)

func NewServiceObject(component *v3.Component, app *v3.Application) corev1.Service {
	ownerRef := GetOwnerRef(app)
	serverPort := component.OptTraits.Ingress.ServerPort
	
	port := corev1.ServicePort {
		Name: "http" + "-" + component.Name,
		Port: serverPort,
		TargetPort: intstr.FromInt(int(serverPort)),
		Protocol: corev1.ProtocolTCP,
	}
	
	service := corev1.Service {
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + "service",
			Annotations:     map[string]string{},
		},
		Spec: corev1.ServiceSpec {
			Selector: map[string]string {
				"app": app.Name + "-" + component.Name + "-" + "workload",
			},
			Ports: []corev1.ServicePort{port},
		},
	}
	
	return service
}

func NewVirtualServiceObject(component *v3.Component, app *v3.Application) istiov1alpha3.VirtualService {
	ownerRef := GetOwnerRef(app)
	host := component.OptTraits.Ingress.Host
	service := app.Name + "-" + component.Name + "-" + "service" + "." + app.Namespace + ".svc.cluster.local"
	port := uint32(component.OptTraits.Ingress.ServerPort)
	
	virtualService := istiov1alpha3.VirtualService {
		TypeMeta: metav1.TypeMeta{
			Kind: "VirtualService",
			APIVersion: "networking.istio.io/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + "vs",
			Annotations:     map[string]string{},
		},
		Spec: istiov1alpha3.VirtualServiceSpec {
			Gateways: []string{(app.Namespace + "-" + "gateway")},
			Hosts: []string{host},
			HTTP:  []istiov1alpha3.HTTPRoute {
				istiov1alpha3.HTTPRoute{
					Route: []istiov1alpha3.HTTPRouteDestination {
						istiov1alpha3.HTTPRouteDestination {
							Destination: istiov1alpha3.Destination{
								Host: service,
								Port: istiov1alpha3.PortSelector{
									Number: port,
								},
							},
						},
					},
				},
			},
		},
	}
	
	return virtualService
}

func NewDestinationruleObject(component *v3.Component, app *v3.Application) istiov1alpha3.DestinationRule {
	ownerRef := GetOwnerRef(app)
	service := app.Name + "-" + component.Name + "-" + "service" + "." + app.Namespace + ".svc.cluster.local"
	
	var lbSetting  *istiov1alpha3.LoadBalancerSettings
	if component.DevTraits.IngressLB.ConsistentType != "" {
		lbSetting = &istiov1alpha3.LoadBalancerSettings {
			ConsistentHash: &istiov1alpha3.ConsistentHashLB {
				UseSourceIP: true,
			},
		}	
	}else if lbType := component.DevTraits.IngressLB.LBType; lbType != "" {
		var simplb  istiov1alpha3.SimpleLB
		switch lbType {
			case "rr":
			   simplb = istiov1alpha3.SimpleLBRoundRobin
			case "leastConn":
			   simplb = istiov1alpha3.SimpleLBLeastConn
			case "random":
			   simplb = istiov1alpha3.SimpleLBRandom
		}
		
		lbSetting = &istiov1alpha3.LoadBalancerSettings {
			Simple: simplb,
		}
	}
	
	destinationrule := istiov1alpha3.DestinationRule {
		TypeMeta: metav1.TypeMeta{
			Kind: "DestinationRule",
			APIVersion: "networking.istio.io/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + "destinationrule",
			Annotations:     map[string]string{},
		},
		Spec: istiov1alpha3.DestinationRuleSpec {
			Host: service,
			TrafficPolicy: &istiov1alpha3.TrafficPolicy {
				LoadBalancer: lbSetting,
			},
		},
	}
	
	return destinationrule
}