package controller

import (
	v3 "github.com/hd-Li/types/apis/project.cattle.io/v3"
	"github.com/knative/pkg/apis/istio/common/v1alpha1"
	istiov1alpha3 "github.com/knative/pkg/apis/istio/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// NewServiceObject Use for generate ServiceObject
func NewServiceObject(app *v3.Application) corev1.Service {
	//ownerRef := GetOwnerRef(app)
	port := corev1.ServicePort{
		Name:       "http" + "-" + app.Name,
		Port:       app.Spec.OptTraits.Ingress.ServerPort,
		TargetPort: intstr.FromInt(int(app.Spec.OptTraits.Ingress.ServerPort)),
		Protocol:   corev1.ProtocolTCP,
	}

	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(app, v3.SchemeGroupVersion.WithKind("Application"))},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + "service",
			Annotations:     map[string]string{},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":    app.Name + "-" + "workload",
				"inpool": "yes",
			},
			Ports: []corev1.ServicePort{port},
		},
	}

	return service
}

// NewVirtualServiceObject Use for generate VirtualServiceObject
func NewVirtualServiceObject(app *v3.Application) istiov1alpha3.VirtualService {
	host := app.Spec.OptTraits.Ingress.Host
	service := app.Name + "-" + "service" + "." + app.Namespace + ".svc.cluster.local"
	port := uint32(app.Spec.OptTraits.Ingress.ServerPort)
	//var matchlist []istiov1alpha3.HTTPMatchRequest

	var httproutes []istiov1alpha3.HTTPRoute
	var httproute istiov1alpha3.HTTPRoute
	httproute = istiov1alpha3.HTTPRoute{
		Match: []istiov1alpha3.HTTPMatchRequest{
			{
				Uri: &v1alpha1.StringMatch{
					Prefix: app.Spec.OptTraits.Ingress.Path,
				},
			},
		},
		/*		Route: []istiov1alpha3.DestinationWeight{
				{
					Destination: istiov1alpha3.Destination{
						Host: service,
						Port: istiov1alpha3.PortSelector{
							Number: port,
						},
					},
				},
			},*/
	}
	// add GrayRelease handlelogic
	if len(app.Spec.Components) < 2 {
		app.Spec.OptTraits.GrayRelease = nil
		httproute.Route = []istiov1alpha3.DestinationWeight{
			{
				Destination: istiov1alpha3.Destination{
					Host: service,
					Port: istiov1alpha3.PortSelector{
						Number: port,
					},
				},
			},
		}
	} else if len(app.Spec.OptTraits.GrayRelease) >= 2 && len(app.Spec.Components) >= 2 {
		for version, weight := range app.Spec.OptTraits.GrayRelease {
			httproute.Route = append(httproute.Route, istiov1alpha3.DestinationWeight{
				Destination: istiov1alpha3.Destination{
					Host: service,
					Port: istiov1alpha3.PortSelector{
						Number: port,
					},
					Subset: version,
				},
				Weight: weight,
			})
		}
	}
	//if !(reflect.DeepEqual(app.Spec.OptTraits.HTTPRetry, v3.HTTPRetry{})) {
	if app.Spec.OptTraits.HTTPRetry != nil {
		httproute.Retries = &istiov1alpha3.HTTPRetry{
			Attempts:      app.Spec.OptTraits.HTTPRetry.Attempts,
			PerTryTimeout: app.Spec.OptTraits.HTTPRetry.PerTryTimeout,
			RetryOn:       "5xx,gateway-error,connect-failure,refused-stream",
		}
	} else {
		httproute.Retries = &istiov1alpha3.HTTPRetry{
			Attempts:      3,
			PerTryTimeout: "10s",
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
			Name:            app.Name + "-" + "vs",
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
func NewDestinationruleObject(app *v3.Application) istiov1alpha3.DestinationRule {
	service := app.Name + "-" + "service" + "." + app.Namespace + ".svc.cluster.local"
	trafficPolicy := new(istiov1alpha3.TrafficPolicy)
	destinationrule := istiov1alpha3.DestinationRule{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DestinationRule",
			APIVersion: "networking.istio.io/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(app, v3.SchemeGroupVersion.WithKind("Application"))},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + "destinationrule",
			Annotations:     map[string]string{},
		},
		Spec: istiov1alpha3.DestinationRuleSpec{
			Host:          service,
			TrafficPolicy: trafficPolicy,
		},
	}

	if len(app.Spec.OptTraits.GrayRelease) >= 2 {
		for k := range app.Spec.OptTraits.GrayRelease {
			var labels map[string]string = make(map[string]string)
			labels["version"] = k
			labels["app"] = app.Name + "-" + "workload"
			labels["inpool"] = "yes"
			destinationrule.Spec.Subsets = append(destinationrule.Spec.Subsets, istiov1alpha3.Subset{
				Name:   k,
				Labels: labels,
			})
		}
	}
	if app.Spec.OptTraits.LoadBalancer != nil {
		if app.Spec.OptTraits.LoadBalancer.ConsistentHash != nil {
			if app.Spec.OptTraits.LoadBalancer.ConsistentHash.UseSourceIP {
				lbsetting := new(istiov1alpha3.LoadBalancerSettings)
				hashlb := new(istiov1alpha3.ConsistentHashLB)
				hashlb = &istiov1alpha3.ConsistentHashLB{
					UseSourceIp: true,
				}
				lbsetting.ConsistentHash = hashlb
				trafficPolicy.LoadBalancer = lbsetting
			}
		} else if lbType := app.Spec.OptTraits.LoadBalancer.Simple; lbType != "" {
			lbsetting := new(istiov1alpha3.LoadBalancerSettings)
			switch lbType {
			case "ROUND_ROBIN":
				lbsetting.Simple = istiov1alpha3.SimpleLBRoundRobin
			case "LEAST_CONN":
				lbsetting.Simple = istiov1alpha3.SimpleLBLeastConn
			case "RANDOM":
				lbsetting.Simple = istiov1alpha3.SimpleLBRandom
			}
			trafficPolicy.LoadBalancer = lbsetting
		}
	} else {
		trafficPolicy.LoadBalancer = &istiov1alpha3.LoadBalancerSettings{Simple: istiov1alpha3.SimpleLBRoundRobin}
	}
	if app.Spec.OptTraits.CircuitBreaking != nil {
		if app.Spec.OptTraits.CircuitBreaking.ConnectionPool != nil {
			trafficPolicy.ConnectionPool = new(istiov1alpha3.ConnectionPoolSettings)
			if app.Spec.OptTraits.CircuitBreaking.ConnectionPool.TCP != nil {
				trafficPolicy.ConnectionPool.Tcp = &istiov1alpha3.TCPSettings{
					MaxConnections: app.Spec.OptTraits.CircuitBreaking.ConnectionPool.TCP.MaxConnections,
					ConnectTimeout: app.Spec.OptTraits.CircuitBreaking.ConnectionPool.TCP.ConnectTimeout,
				}
			}
			if app.Spec.OptTraits.CircuitBreaking.ConnectionPool.HTTP != nil {
				trafficPolicy.ConnectionPool.Http = &istiov1alpha3.HTTPSettings{
					Http1MaxPendingRequests:  app.Spec.OptTraits.CircuitBreaking.ConnectionPool.HTTP.HTTP1MaxPendingRequests,
					Http2MaxRequests:         app.Spec.OptTraits.CircuitBreaking.ConnectionPool.HTTP.HTTP2MaxRequests,
					MaxRequestsPerConnection: app.Spec.OptTraits.CircuitBreaking.ConnectionPool.HTTP.MaxRequestsPerConnection,
					MaxRetries:               app.Spec.OptTraits.CircuitBreaking.ConnectionPool.HTTP.MaxRetries,
				}
			}
		}
		if app.Spec.OptTraits.CircuitBreaking.OutlierDetection != nil {
			trafficPolicy.OutlierDetection = &istiov1alpha3.OutlierDetection{
				ConsecutiveErrors:  app.Spec.OptTraits.CircuitBreaking.OutlierDetection.ConsecutiveErrors,
				Interval:           app.Spec.OptTraits.CircuitBreaking.OutlierDetection.Interval,
				BaseEjectionTime:   app.Spec.OptTraits.CircuitBreaking.OutlierDetection.BaseEjectionTime,
				MaxEjectionPercent: app.Spec.OptTraits.CircuitBreaking.OutlierDetection.MaxEjectionPercent,
			}
		}
	}
	return destinationrule
}
