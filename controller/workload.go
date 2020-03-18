package controller

import (
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"

	v3 "github.com/hd-Li/types/apis/project.cattle.io/v3"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/credentialprovider"
)

func NewConfigMapObject(component *v3.Component, app *v3.Application) corev1.ConfigMap {
	ownerRef := GetOwnerRef(app)
	var stringmap map[string]string = make(map[string]string)
	for _, i := range component.Containers {
		for _, j := range i.Config {
			stringmap[j.FileName] = j.Value
		}
	}
	configmap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + component.Version + "-" + "configmap",
		},
		Data: stringmap,
	}
	return configmap
}

func NewSecretObject(component *v3.Component, app *v3.Application) corev1.Secret {
	ownerRef := GetOwnerRef(app)

	dockercfgJsonContent, err := handleDockerCfgJsonContent(component.DevTraits.ImagePullConfig.Username, component.DevTraits.ImagePullConfig.Password, "", component.DevTraits.ImagePullConfig.Registry)
	if err != nil {
		log.Errorf("Create docker secret failed for %s %s ", app.Namespace, component.Name)
		return corev1.Secret{}
	}
	datamap := map[string][]byte{}
	datamap[corev1.DockerConfigJsonKey] = dockercfgJsonContent
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + "registry-secret",
		},
		Data: datamap,
		Type: corev1.SecretTypeDockerConfigJson,
	}
	return secret
}

// handleDockerCfgJsonContent serializes a ~/.docker/config.json file
func handleDockerCfgJsonContent(username, password, email, server string) ([]byte, error) {
	dockercfgAuth := credentialprovider.DockerConfigEntry{
		Username: username,
		Password: password,
		Email:    email,
	}

	dockerCfgJson := credentialprovider.DockerConfigJson{
		Auths: map[string]credentialprovider.DockerConfigEntry{server: dockercfgAuth},
	}

	return json.Marshal(dockerCfgJson)
}

func NewDeployObject(component *v3.Component, app *v3.Application) appsv1beta2.Deployment {
	ownerRef := GetOwnerRef(app)
	var volumes []corev1.Volume //zk
	for _, i := range component.Containers {
		for _, j := range i.Resources.Volumes {
			if j.Disk.Ephemeral {
				volumes = append(volumes, corev1.Volume{Name: component.Name + j.Name,
					VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
				})
			} else {
				var pathtype corev1.HostPathType = corev1.HostPathDirectoryOrCreate
				volumes = append(volumes, corev1.Volume{Name: component.Name + j.Name,
					VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: j.Disk.Required,
						Type: &pathtype},
					}})
			}
		}
		for _, k := range i.Config {
			volumes = append(volumes, corev1.Volume{Name: component.Name + "-" + component.Version + "-" + strings.Replace(k.FileName, ".", "-", -1),
				VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: app.Name + "-" + component.Name + component.Version + "-" + "configmap"},
					Items: []corev1.KeyToPath{corev1.KeyToPath{
						Key:  k.FileName,
						Path: "tmp/" + k.FileName,
					}},
				}}})
		}
	}
	containers, _ := getContainers(component)
	var imagepullsecret []corev1.LocalObjectReference
	if app.Status.ComponentResource[(app.Name+"_"+component.Name+"_"+component.Version)].ImagePullSecret != "" {
		imagepullsecret = append(imagepullsecret, corev1.LocalObjectReference{app.Status.ComponentResource[(app.Name + "_" + component.Name + "_" + component.Version)].ImagePullSecret})
	}
	deploy := appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + "workload" + "-" + component.Version,
			Labels:          app.Labels,
			Annotations:     app.Annotations,
		},
		Spec: appsv1beta2.DeploymentSpec{
			//add replicas
			Replicas: &component.OptTraits.ManualScaler.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":     app.Name + "-" + component.Name + "-" + "workload",
					"version": component.Version,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":     app.Name + "-" + component.Name + "-" + "workload",
						"version": component.Version,
						"inpool":  "yes",
					},
				},

				Spec: corev1.PodSpec{
					ImagePullSecrets: imagepullsecret,
					Containers:       containers,
					Volumes:          volumes, //zk
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
		livenesshandler, readinesshandler := getContainersHealthCheck(cc)
		var volumes []corev1.VolumeMount
		for _, j := range cc.Resources.Volumes {
			volumes = append(volumes, corev1.VolumeMount{
				Name:      component.Name + j.Name,
				MountPath: j.MountPath,
			})
		}
		for _, k := range cc.Config {
			volumes = append(volumes, corev1.VolumeMount{
				Name:      component.Name + "-" + component.Version + "-" + strings.Replace(k.FileName, ".", "-", -1),
				MountPath: k.Path + "/" + k.FileName,
				SubPath:   "tmp/" + k.FileName,
			})
		}

		container := corev1.Container{
			Name:         cc.Name,
			Image:        cc.Image,
			Command:      cc.Command,
			Args:         cc.Args,
			Ports:        ports,
			Env:          envs,
			Resources:    resources,
			VolumeMounts: volumes,
			LivenessProbe: &corev1.Probe{
				InitialDelaySeconds: cc.LivenessProbe.InitialDelaySeconds,
				TimeoutSeconds:      cc.LivenessProbe.TimeoutSeconds,
				FailureThreshold:    cc.LivenessProbe.FailureThreshold,
				Handler:             livenesshandler,
			},
			ReadinessProbe: &corev1.Probe{
				InitialDelaySeconds: cc.ReadinessProbe.InitialDelaySeconds,
				PeriodSeconds:       cc.ReadinessProbe.PeriodSeconds,
				TimeoutSeconds:      cc.ReadinessProbe.TimeoutSeconds,
				Handler:             readinesshandler,
			},
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
		corev1.ResourceCPU:    resource.MustParse(cpu),
		corev1.ResourceMemory: resource.MustParse(mem),
	}

	rr := corev1.ResourceRequirements{
		Requests: resources,
		Limits:   resources,
	}

	return rr
}

func getContainerEnvs(cc v3.ComponentContainer) []corev1.EnvVar {
	var envs []corev1.EnvVar

	for _, ccenv := range cc.Env {
		env := corev1.EnvVar{
			Name:  ccenv.Name,
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

		if ccp.Protocol == "tcp" || ccp.Protocol == "" {
			proto = corev1.ProtocolTCP
		} else {
			proto = corev1.ProtocolUDP
		}

		port := corev1.ContainerPort{
			Name:          ccp.Name,
			ContainerPort: ccp.ContainerPort,
			Protocol:      proto,
		}

		ports = append(ports, port)
	}

	return ports
}

// zk generate health check model data
func getContainersHealthCheck(cc v3.ComponentContainer) (livenesshandler corev1.Handler, readinesshandler corev1.Handler) {

	if len(cc.LivenessProbe.Exec.Command) != 0 {
		livenesshandler = corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: cc.ReadinessProbe.Exec.Command,
			},
		}
	} else if cc.LivenessProbe.HTTPGet.Path != "" {
		livenesshandler = corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: cc.ReadinessProbe.HTTPGet.Path,
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: int32(cc.ReadinessProbe.HTTPGet.Port),
				},
			},
		}
	} else {
		livenesshandler = corev1.Handler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					IntVal: int32(cc.ReadinessProbe.TCPSocket.Port),
					Type:   intstr.Int,
				},
			},
		}
	}

	if len(cc.ReadinessProbe.Exec.Command) != 0 {
		readinesshandler = corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: cc.ReadinessProbe.Exec.Command,
			},
		}
	} else if cc.ReadinessProbe.HTTPGet.Path != "" {
		readinesshandler = corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: cc.ReadinessProbe.HTTPGet.Path,
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: int32(cc.ReadinessProbe.HTTPGet.Port),
				},
			},
		}
	} else {
		readinesshandler = corev1.Handler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					IntVal: int32(cc.ReadinessProbe.TCPSocket.Port),
					Type:   intstr.Int,
				},
			},
		}
	}
	return
}
