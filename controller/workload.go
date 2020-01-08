package controller

import (
	"github.com/hd-Li/types/apis/project.cattle.io/v3"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func NewDeployObject(component *v3.Component, app *v3.Application) appsv1beta2.Deployment {
	ownerRef := GetOwnerRef(app)
	containers, _ := getContainers(component)
	deploy := appsv1beta2.Deployment {
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + "workload" + "-" + component.Version,
			Labels:          app.Labels,
			Annotations:     app.Annotations,
		},
		Spec: appsv1beta2.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": app.Name + "-" + component.Name + "-" + "workload",
					"version": component.Version,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": app.Name + "-" + component.Name + "-" + "workload",
						"version": component.Version,
					},
				},
				Spec: corev1.PodSpec{
					Containers: containers,
				},
			},
		},
	}
	
	return deploy
}

func getContainers(component *v3.Component) ([]corev1.Container, error) {
	var containers []corev1.Container
	
	for _, cc := range component.Containers {
		ports := getContainerPorts(cc)
		envs := getContainerEnvs(cc)
		resources := getContainerResources(cc)
		
		container := corev1.Container {
			Name: cc.Name,
			Image: cc.Image,
			Command: cc.Command,
			Args: cc.Args,
			Ports: ports,
			Env: envs,
			Resources: resources,
		}
		
		containers = append(containers, container)
	}
	
	return containers, nil
}

func getContainerResources(cc v3.ComponentContainer) corev1.ResourceRequirements {
	cpu := "500m"
	mem := "200Mi"
	
	if cc.Resources.Cpu != "" {
		cpu = cc.Resources.Cpu
	}
	
	if cc.Resources.Memory != "" {
		mem = cc.Resources.Memory
	}
	
	resources := map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU: resource.MustParse(cpu),
		corev1.ResourceMemory: resource.MustParse(mem),
	}
	
	rr := corev1.ResourceRequirements {
		Requests: resources,
		Limits: resources,
	}
	
	return rr
}

func getContainerEnvs(cc v3.ComponentContainer) []corev1.EnvVar {
	var envs []corev1.EnvVar
	
	for _, ccenv := range cc.Env {
		env := corev1.EnvVar {
			Name: ccenv.Name,
			Value: ccenv.Value,
		}
		
		envs = append(envs, env)
	}
	
	return envs
}

func getContainerPorts(cc v3.ComponentContainer) []corev1.ContainerPort {
	var ports []corev1.ContainerPort
	
	for _, ccp := range cc.Ports {
		var proto corev1.Protocol
		
		if ccp.Protocol == "tcp" || ccp.Protocol == ""{
			proto = corev1.ProtocolTCP
		}else {
			proto = corev1.ProtocolUDP
		}
		
		port := corev1.ContainerPort {
			Name: ccp.Name,
			ContainerPort: ccp.ContainerPort,
			Protocol: proto,
		}
		
		ports = append(ports, port)
	}
	
	return ports
}