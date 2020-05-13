package controller

import (
	"reflect"

	v3 "github.com/hd-Li/types/apis/project.cattle.io/v3"
	"github.com/knative/pkg/apis/istio/common/v1alpha1"
	istiov1alpha3 "github.com/knative/pkg/apis/istio/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// NewServiceObject Use for generate ServiceObject
func NewServiceObject(component *v3.Component, app *v3.Application) corev1.Service {
	//ownerRef := GetOwnerRef(app)
	serverPort := component.OptTraits.Ingress.ServerPort

	port := corev1.ServicePort{
		Name:       "http" + "-" + component.Name,
		Port:       serverPort,
		TargetPort: intstr.FromInt(int(serverPort)),
		Protocol:   corev1.ProtocolTCP,
	}

	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(app, v3.SchemeGroupVersion.WithKind("Application"))},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + "service",
			Annotations:     map[string]string{},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":    app.Name + "-" + component.Name + "-" + "workload",
				"inpool": "yes",
			},
			Ports: []corev1.ServicePort{port},
		},
	}

	return service
}

// NewVirtualServiceObject Use for generate VirtualServiceObject
func NewVirtualServiceObject(component *v3.Component, app *v3.Application) istiov1alpha3.VirtualService {
	host := component.OptTraits.Ingress.Host
	service := app.Name + "-" + component.Name + "-" + "service" + "." + app.Namespace + ".svc.cluster.local"
	port := uint32(component.OptTraits.Ingress.ServerPort)
	//var matchlist []istiov1alpha3.HTTPMatchRequest

	var httproutes []istiov1alpha3.HTTPRoute
	var httproute istiov1alpha3.HTTPRoute
	httproute = istiov1alpha3.HTTPRoute{
		Match: []istiov1alpha3.HTTPMatchRequest{
			{
				Uri: &v1alpha1.StringMatch{
					Prefix: component.OptTraits.Ingress.Path,
				},
			},
		},
		Route: []istiov1alpha3.DestinationWeight{
			{
				Destination: istiov1alpha3.Destination{
					Host: service,
					Port: istiov1alpha3.PortSelector{
						Number: port,
					},
				},
			},
		},
	}

	if !(reflect.DeepEqual(component.OptTraits.HttpRetry, v3.HttpRetry{})) {
		httproute.Retries = &istiov1alpha3.HTTPRetry{
			Attempts:      component.OptTraits.HttpRetry.Attempts,
			PerTryTimeout: component.OptTraits.HttpRetry.PerTryTimeout,
			RetryOn:       "5xx,gateway-error,connect-failure,refused-stream",
		}
	}

	httproutes = append(httproutes, httproute)

	virtualService := istiov1alpha3.VirtualService{
		TypeMeta: metav1.TypeMeta{
			Kind:       "VirtualService",
			APIVersion: "networking.istio.io/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(app, v3.SchemeGroupVersion.WithKind("Application"))},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + "vs",
			Annotations:     map[string]string{},
		},
		Spec: istiov1alpha3.VirtualServiceSpec{
			Gateways: []string{(app.Namespace + "-" + "gateway")},
			Hosts:    []string{host},
			Http:     httproutes,
		},
	}

	return virtualService
}

// NewDestinationruleObject Use for generate DestinationruleObject
func NewDestinationruleObject(component *v3.Component, app *v3.Application) istiov1alpha3.DestinationRule {
	service := app.Name + "-" + component.Name + "-" + "service" + "." + app.Namespace + ".svc.cluster.local"

	var lbSetting *istiov1alpha3.LoadBalancerSettings
	var connectionpoolsetting *istiov1alpha3.ConnectionPoolSettings //zk
	var outlierdetectionsetting *istiov1alpha3.OutlierDetection     //zk
	if component.DevTraits.IngressLB.ConsistentType != "" {
		lbSetting = &istiov1alpha3.LoadBalancerSettings{
			ConsistentHash: &istiov1alpha3.ConsistentHashLB{
				UseSourceIp: true,
			},
		}
	} else if lbType := component.DevTraits.IngressLB.LBType; lbType != "" {
		var simplb istiov1alpha3.SimpleLB
		switch lbType {
		case "rr":
			simplb = istiov1alpha3.SimpleLBRoundRobin
		case "leastConn":
			simplb = istiov1alpha3.SimpleLBLeastConn
		case "random":
			simplb = istiov1alpha3.SimpleLBRandom
		}

		lbSetting = &istiov1alpha3.LoadBalancerSettings{
			Simple: simplb,
		}
	}
	connectionpoolsetting = &istiov1alpha3.ConnectionPoolSettings{
		Tcp: &istiov1alpha3.TCPSettings{
			MaxConnections: component.OptTraits.CircuitBreaking.ConnectionPool.TCP.MaxConnections,
			ConnectTimeout: component.OptTraits.CircuitBreaking.ConnectionPool.TCP.ConnectTimeout,
		},
		Http: &istiov1alpha3.HTTPSettings{
			Http1MaxPendingRequests:  component.OptTraits.CircuitBreaking.ConnectionPool.HTTP.HTTP1MaxPendingRequests,
			Http2MaxRequests:         component.OptTraits.CircuitBreaking.ConnectionPool.HTTP.HTTP2MaxRequests,
			MaxRequestsPerConnection: component.OptTraits.CircuitBreaking.ConnectionPool.HTTP.MaxRequestsPerConnection,
			MaxRetries:               component.OptTraits.CircuitBreaking.ConnectionPool.HTTP.MaxRetries,
		},
	}
	outlierdetectionsetting = &istiov1alpha3.OutlierDetection{
		ConsecutiveErrors:  component.OptTraits.CircuitBreaking.OutlierDetection.ConsecutiveErrors,
		Interval:           component.OptTraits.CircuitBreaking.OutlierDetection.Interval,
		BaseEjectionTime:   component.OptTraits.CircuitBreaking.OutlierDetection.BaseEjectionTime,
		MaxEjectionPercent: component.OptTraits.CircuitBreaking.OutlierDetection.MaxEjectionPercent,
	}
	destinationrule := istiov1alpha3.DestinationRule{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DestinationRule",
			APIVersion: "networking.istio.io/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(app, v3.SchemeGroupVersion.WithKind("Application"))},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + "destinationrule",
			Annotations:     map[string]string{},
		},
		Spec: istiov1alpha3.DestinationRuleSpec{
			Host: service,
			TrafficPolicy: &istiov1alpha3.TrafficPolicy{
				LoadBalancer:     lbSetting,
				ConnectionPool:   connectionpoolsetting,
				OutlierDetection: outlierdetectionsetting,
			},
		},
	}
	return destinationrule
}
