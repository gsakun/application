package controller

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"reflect"

	v3 "github.com/hd-Li/types/apis/project.cattle.io/v3"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/credentialprovider"
)

// NewConfigMapObject Use for generate ConfigMapObject
func NewConfigMapObject(component *v3.Component, app *v3.Application) corev1.ConfigMap {
	var stringmap map[string]string = make(map[string]string)
	for _, i := range component.Containers {
		for _, j := range i.Config {
			if j.FileName == "" {
				log.Errorf("%s-%s's configmap configuration's filename is nil,please check configration")
				continue
			}
			stringmap[j.FileName] = j.Value
		}
	}
	configmap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(app, v3.SchemeGroupVersion.WithKind("Application"))},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + component.Version + "-" + "configmap",
		},
		Data: stringmap,
	}
	return configmap
}

// NewSecretObject Use for generate SecretObject
func NewSecretObject(component *v3.Component, app *v3.Application) corev1.Secret {
	dockercfgJSONContent, err := handleDockerCfgJSONContent(component.DevTraits.ImagePullConfig.Username, component.DevTraits.ImagePullConfig.Password, "", component.DevTraits.ImagePullConfig.Registry)
	if err != nil {
		log.Errorf("Create docker secret failed for %s %s ", app.Namespace, component.Name)
		return corev1.Secret{}
	}
	datamap := map[string][]byte{}
	datamap[corev1.DockerConfigJsonKey] = dockercfgJSONContent
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(app, v3.SchemeGroupVersion.WithKind("Application"))},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + "registry-secret",
		},
		Data: datamap,
		Type: corev1.SecretTypeDockerConfigJson,
	}
	return secret
}

// handleDockerCfgJSONContent serializes a ~/.docker/config.json file
func handleDockerCfgJSONContent(username, password, email, server string) ([]byte, error) {
	dockercfgAuth := credentialprovider.DockerConfigEntry{
		Username: username,
		Password: password,
		Email:    email,
	}

	dockerCfgJSON := credentialprovider.DockerConfigJson{
		Auths: map[string]credentialprovider.DockerConfigEntry{server: dockercfgAuth},
	}

	return json.Marshal(dockerCfgJSON)
}

// NewDeployObject Use for generate DeployObject
func NewDeployObject(component *v3.Component, app *v3.Application) appsv1beta2.Deployment {
	//ownerRef := GetOwnerRef(app)
	var volumes []corev1.Volume //zk
	for _, i := range component.Containers {
		for _, j := range i.Resources.Volumes {
			if j.Disk.Ephemeral {
				volumes = append(volumes, corev1.Volume{Name: component.Name + "-" + j.Name,
					VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
				})
			} else {
				var pathtype corev1.HostPathType = corev1.HostPathDirectoryOrCreate
				volumes = append(volumes, corev1.Volume{Name: component.Name + "-" + j.Name,
					VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: j.Disk.Required,
						Type: &pathtype},
					}})
			}
		}
		for _, k := range i.Config {
			if k.FileName == "" {
				continue
			}
			volumes = append(volumes, corev1.Volume{Name: component.Name + "-" + component.Version + "-" + strings.Replace(k.FileName, ".", "-", -1),
				VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: app.Name + "-" + component.Name + "-" + component.Version + "-" + "configmap"},
					Items: []corev1.KeyToPath{
						{
							Key:  k.FileName,
							Path: "tmp/" + k.FileName,
						}},
				}}})
		}
	}
	containers, _ := getContainers(component)
	var imagepullsecret []corev1.LocalObjectReference
	if app.Status.ComponentResource[(app.Name+"_"+component.Name+"_"+component.Version)].ImagePullSecret != "" {
		imagepullsecret = append(imagepullsecret, corev1.LocalObjectReference{Name: app.Status.ComponentResource[(app.Name + "_" + component.Name + "_" + component.Version)].ImagePullSecret})
	}
	deploy := appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(app, v3.SchemeGroupVersion.WithKind("Application"))},
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
	if !reflect.DeepEqual(component.OptTraits.SchedulePolicy, v3.SchedulePolicy{}) {
		if component.OptTraits.SchedulePolicy.NodeSelector != nil {
			deploy.Spec.Template.Spec.NodeSelector = component.OptTraits.SchedulePolicy.NodeSelector
		}
		if !reflect.DeepEqual(component.OptTraits.SchedulePolicy.NodeAffinity, v3.CNodeAffinity{}) {
			deploy.Spec.Template.Spec.Affinity = new(corev1.Affinity)
			if component.OptTraits.SchedulePolicy.NodeAffinity.HardAffinity {
				deploy.Spec.Template.Spec.Affinity.NodeAffinity = new(corev1.NodeAffinity)
				deploy.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = &corev1.NodeSelector{
					NodeSelectorTerms: []corev1.NodeSelectorTerm{
						{
							MatchExpressions: []corev1.NodeSelectorRequirement{
								{
									Key:      component.OptTraits.SchedulePolicy.NodeAffinity.CLabelSelectorRequirement.Key,
									Operator: corev1.NodeSelectorOperator(component.OptTraits.SchedulePolicy.NodeAffinity.CLabelSelectorRequirement.Operator),
									Values:   component.OptTraits.SchedulePolicy.NodeAffinity.CLabelSelectorRequirement.Values,
								},
							},
						},
					},
				}
			} else {
				deploy.Spec.Template.Spec.Affinity.NodeAffinity = new(corev1.NodeAffinity)
				deploy.Spec.Template.Spec.Affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution = []corev1.PreferredSchedulingTerm{
					{
						Weight: int32(1),
						Preference: corev1.NodeSelectorTerm{
							MatchExpressions: []corev1.NodeSelectorRequirement{
								{
									Key:      component.OptTraits.SchedulePolicy.NodeAffinity.CLabelSelectorRequirement.Key,
									Operator: corev1.NodeSelectorOperator(component.OptTraits.SchedulePolicy.NodeAffinity.CLabelSelectorRequirement.Operator),
									Values:   component.OptTraits.SchedulePolicy.NodeAffinity.CLabelSelectorRequirement.Values,
								},
							},
						},
					},
				}
			}
		}
		if !reflect.DeepEqual(component.OptTraits.SchedulePolicy.PodAffinity, v3.CPodAffinity{}) {
			if deploy.Spec.Template.Spec.Affinity != nil {
				deploy.Spec.Template.Spec.Affinity = new(corev1.Affinity)
			}
			if component.OptTraits.SchedulePolicy.PodAffinity.HardAffinity {
				deploy.Spec.Template.Spec.Affinity.PodAffinity = new(corev1.PodAffinity)
				deploy.Spec.Template.Spec.Affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution = []corev1.PodAffinityTerm{
					{
						TopologyKey: "kubernetes.io/hostname",
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      component.OptTraits.SchedulePolicy.PodAffinity.CLabelSelectorRequirement.Key,
									Operator: metav1.LabelSelectorOperator(component.OptTraits.SchedulePolicy.PodAffinity.CLabelSelectorRequirement.Operator),
									Values:   component.OptTraits.SchedulePolicy.PodAffinity.CLabelSelectorRequirement.Values,
								},
							},
						},
					},
				}
			} else {
				deploy.Spec.Template.Spec.Affinity.PodAffinity = new(corev1.PodAffinity)
				deploy.Spec.Template.Spec.Affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution = []corev1.WeightedPodAffinityTerm{
					{
						Weight: int32(1),
						PodAffinityTerm: corev1.PodAffinityTerm{
							TopologyKey: "kubernetes.io/hostname",
							LabelSelector: &metav1.LabelSelector{
								MatchExpressions: []metav1.LabelSelectorRequirement{
									{
										Key:      component.OptTraits.SchedulePolicy.PodAffinity.CLabelSelectorRequirement.Key,
										Operator: metav1.LabelSelectorOperator(component.OptTraits.SchedulePolicy.PodAffinity.CLabelSelectorRequirement.Operator),
										Values:   component.OptTraits.SchedulePolicy.PodAffinity.CLabelSelectorRequirement.Values,
									},
								},
							},
						},
					},
				}
			}
		}

		if !reflect.DeepEqual(component.OptTraits.SchedulePolicy.PodAntiAffinity, v3.CPodAntiAffinity{}) {
			if deploy.Spec.Template.Spec.Affinity != nil {
				deploy.Spec.Template.Spec.Affinity = new(corev1.Affinity)
			}
			if component.OptTraits.SchedulePolicy.PodAntiAffinity.HardAffinity {
				deploy.Spec.Template.Spec.Affinity.PodAntiAffinity = new(corev1.PodAntiAffinity)
				deploy.Spec.Template.Spec.Affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = []corev1.PodAffinityTerm{
					{
						TopologyKey: "kubernetes.io/hostname",
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      component.OptTraits.SchedulePolicy.PodAntiAffinity.CLabelSelectorRequirement.Key,
									Operator: metav1.LabelSelectorOperator(component.OptTraits.SchedulePolicy.PodAntiAffinity.CLabelSelectorRequirement.Operator),
									Values:   component.OptTraits.SchedulePolicy.PodAntiAffinity.CLabelSelectorRequirement.Values,
								},
							},
						},
					},
				}
			} else {
				deploy.Spec.Template.Spec.Affinity.PodAntiAffinity = new(corev1.PodAntiAffinity)
				deploy.Spec.Template.Spec.Affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution = []corev1.WeightedPodAffinityTerm{
					{
						Weight: int32(1),
						PodAffinityTerm: corev1.PodAffinityTerm{
							TopologyKey: "kubernetes.io/hostname",
							LabelSelector: &metav1.LabelSelector{
								MatchExpressions: []metav1.LabelSelectorRequirement{
									{
										Key:      component.OptTraits.SchedulePolicy.PodAntiAffinity.CLabelSelectorRequirement.Key,
										Operator: metav1.LabelSelectorOperator(component.OptTraits.SchedulePolicy.PodAntiAffinity.CLabelSelectorRequirement.Operator),
										Values:   component.OptTraits.SchedulePolicy.PodAntiAffinity.CLabelSelectorRequirement.Values,
									},
								},
							},
						},
					},
				}
			}
		}

	}
	if component.OptTraits.TerminationGracePeriodSeconds > 30 {
		deploy.Spec.Template.Spec.TerminationGracePeriodSeconds = &component.OptTraits.TerminationGracePeriodSeconds
	}
	if component.OptTraits.CustomMetric.Enable && component.OptTraits.CustomMetric.Uri != "" {
		deploy.Spec.Template.Annotations = make(map[string]string)
		deploy.Spec.Template.Annotations["prometheus.io/path"] = "/metrics"
		deploy.Spec.Template.Annotations["prometheus.io/port"] = "16666"
		deploy.Spec.Template.Annotations["prometheus.io/scrape"] = "true"
	}
	// TODO OPEN METRIC TRANSFR
	return deploy
}

func getContainers(component *v3.Component) ([]corev1.Container, error) {
	var containers []corev1.Container
	for _, cc := range component.Containers {
		ports := getContainerPorts(cc)
		envs := getContainerEnvs(cc)
		resources := getContainerResources(cc)
		livenesshandler, readinesshandler := getContainersHealthCheck(cc)
		lifecycle := getContainersLifeCycle(cc)
		var volumes []corev1.VolumeMount
		for _, j := range cc.Resources.Volumes {
			volumes = append(volumes, corev1.VolumeMount{
				Name:      component.Name + "-" + j.Name,
				MountPath: j.MountPath,
			})
		}
		for _, k := range cc.Config {
			if k.FileName == "" {
				continue
			}
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
		}
		if lifecycle != nil {
			container.Lifecycle = lifecycle
		}
		if !(reflect.DeepEqual(livenesshandler, corev1.Handler{})) {
			container.LivenessProbe = &corev1.Probe{
				InitialDelaySeconds: cc.LivenessProbe.InitialDelaySeconds,
				TimeoutSeconds:      cc.LivenessProbe.TimeoutSeconds,
				FailureThreshold:    cc.LivenessProbe.FailureThreshold,
				Handler:             livenesshandler,
			}
		}
		if !(reflect.DeepEqual(readinesshandler, corev1.Handler{})) {
			container.ReadinessProbe = &corev1.Probe{
				InitialDelaySeconds: cc.ReadinessProbe.InitialDelaySeconds,
				PeriodSeconds:       cc.ReadinessProbe.PeriodSeconds,
				TimeoutSeconds:      cc.ReadinessProbe.TimeoutSeconds,
				Handler:             readinesshandler,
			}
		}
		containers = append(containers, container)
		if component.OptTraits.CustomMetric.Enable && component.OptTraits.CustomMetric.Uri != "" {
			containers = append(containers, corev1.Container{
				Name:            "transter-proxy",
				Image:           os.Getenv("PROXYIMAGE"),
				ImagePullPolicy: corev1.PullIfNotPresent,
				Env: []corev1.EnvVar{
					{
						Name: "POD_NAME",
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								APIVersion: "v1",
								FieldPath:  "metadata.name",
							}},
					},
					{
						Name: "POD_NAMESPACE",
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								APIVersion: "v1",
								FieldPath:  "metadata.namespace",
							}},
					},
					{
						Name: "POD_IP",
						ValueFrom: &corev1.EnvVarSource{
							FieldRef: &corev1.ObjectFieldSelector{
								APIVersion: "v1",
								FieldPath:  "status.podIP",
							}},
					},
					{
						Name:  "URI",
						Value: component.OptTraits.CustomMetric.Uri,
					},
				},
			})
		}
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
	if cc.Resources.Gpu > 0 {
		resources[corev1.ResourceName("nvidia.com/gpu")] = resource.MustParse(strconv.Itoa(cc.Resources.Gpu))
	}
	rr := corev1.ResourceRequirements{
		Requests: resources,
		Limits:   resources,
	}

	return rr
}

func getContainerEnvs(cc v3.ComponentContainer) (envs []corev1.EnvVar) {
	for _, ccenv := range cc.Env {
		if ccenv.FromParam != "" && (ccenv.FromParam == "spec.nodeName" || ccenv.FromParam == "metadata.name" || ccenv.FromParam == "metadata.namespace" || ccenv.FromParam == "status.podIP") {
			env := corev1.EnvVar{
				Name: ccenv.Name,
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: ccenv.FromParam,
					}},
			}
			envs = append(envs, env)
		} else {
			env := corev1.EnvVar{
				Name:  ccenv.Name,
				Value: ccenv.Value,
			}
			envs = append(envs, env)
		}

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
	//log.Debugf("Container info is %v", cc)
	if !reflect.DeepEqual(cc.LivenessProbe, v3.HealthProbe{}) {
		if len(cc.LivenessProbe.Exec.Command) != 0 {
			livenesshandler = corev1.Handler{
				Exec: &corev1.ExecAction{
					Command: cc.LivenessProbe.Exec.Command,
				},
			}
		} else if cc.LivenessProbe.HTTPGet.Path != "" && cc.LivenessProbe.HTTPGet.Port > 0 {
			livenesshandler = corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: cc.LivenessProbe.HTTPGet.Path,
					Port: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(cc.LivenessProbe.HTTPGet.Port),
					},
				},
			}
		} else if cc.LivenessProbe.TCPSocket.Port > 0 {
			livenesshandler = corev1.Handler{
				TCPSocket: &corev1.TCPSocketAction{
					Port: intstr.IntOrString{
						IntVal: int32(cc.LivenessProbe.TCPSocket.Port),
						Type:   intstr.Int,
					},
				},
			}
		} else {
			livenesshandler = corev1.Handler{}
		}
	}
	if !reflect.DeepEqual(cc.ReadinessProbe, v3.HealthProbe{}) {
		if len(cc.ReadinessProbe.Exec.Command) != 0 {
			readinesshandler = corev1.Handler{
				Exec: &corev1.ExecAction{
					Command: cc.ReadinessProbe.Exec.Command,
				},
			}
		} else if cc.ReadinessProbe.HTTPGet.Path != "" && cc.ReadinessProbe.HTTPGet.Port > 0 {
			readinesshandler = corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: cc.ReadinessProbe.HTTPGet.Path,
					Port: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(cc.ReadinessProbe.HTTPGet.Port),
					},
				},
			}
		} else if cc.ReadinessProbe.TCPSocket.Port > 0 {
			readinesshandler = corev1.Handler{
				TCPSocket: &corev1.TCPSocketAction{
					Port: intstr.IntOrString{
						IntVal: int32(cc.ReadinessProbe.TCPSocket.Port),
						Type:   intstr.Int,
					},
				},
			}
		} else {
			readinesshandler = corev1.Handler{}
		}
	}
	return
}

// zk generate pod lifecycle
func getContainersLifeCycle(cc v3.ComponentContainer) (lifecycle *corev1.Lifecycle) {
	// new lifecycle memory address
	lifecycle = new(corev1.Lifecycle)
	if reflect.DeepEqual(cc.Lifecycle, v3.CLifecycle{}) {
		return lifecycle
	}
	log.Debugf("Container info is %v", cc)
	if cc.Lifecycle.PostStart != nil {
		if !reflect.DeepEqual(cc.Lifecycle.PostStart.Exec, v3.ExecAction{}) {
			if len(cc.Lifecycle.PostStart.Exec.Command) != 0 {
				lifecycle.PostStart = &corev1.Handler{
					Exec: &corev1.ExecAction{
						Command: cc.Lifecycle.PostStart.Exec.Command,
					},
				}
			}
		}
		if !reflect.DeepEqual(cc.Lifecycle.PostStart.HTTPGet, v3.HTTPGetAction{}) {
			if cc.Lifecycle.PostStart.HTTPGet.Path != "" && cc.Lifecycle.PostStart.HTTPGet.Port > 0 {
				lifecycle.PostStart = &corev1.Handler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: cc.Lifecycle.PostStart.HTTPGet.Path,
						Port: intstr.IntOrString{
							Type:   intstr.Int,
							IntVal: int32(cc.Lifecycle.PostStart.HTTPGet.Port),
						},
					},
				}
			}
		}
		if !reflect.DeepEqual(cc.Lifecycle.PostStart.TCPSocket, v3.TCPSocketAction{}) {
			if cc.Lifecycle.PostStart.TCPSocket.Port > 0 {
				lifecycle.PostStart = &corev1.Handler{
					TCPSocket: &corev1.TCPSocketAction{
						Port: intstr.IntOrString{
							IntVal: int32(cc.Lifecycle.PostStart.TCPSocket.Port),
							Type:   intstr.Int,
						},
					},
				}
			}
		}
	}
	if cc.Lifecycle.PreStop != nil {
		if !reflect.DeepEqual(cc.Lifecycle.PreStop.Exec, v3.ExecAction{}) {
			if len(cc.Lifecycle.PreStop.Exec.Command) != 0 {
				lifecycle.PreStop = &corev1.Handler{
					Exec: &corev1.ExecAction{
						Command: cc.Lifecycle.PreStop.Exec.Command,
					},
				}
			}
		}
		if !reflect.DeepEqual(cc.Lifecycle.PreStop.HTTPGet, v3.HTTPGetAction{}) {
			if cc.Lifecycle.PreStop.HTTPGet.Path != "" && cc.Lifecycle.PreStop.HTTPGet.Port > 0 {
				lifecycle.PreStop = &corev1.Handler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: cc.Lifecycle.PreStop.HTTPGet.Path,
						Port: intstr.IntOrString{
							Type:   intstr.Int,
							IntVal: int32(cc.Lifecycle.PreStop.HTTPGet.Port),
						},
					},
				}
			}
		}
		if !reflect.DeepEqual(cc.Lifecycle.PreStop.TCPSocket, v3.TCPSocketAction{}) {
			if cc.Lifecycle.PreStop.TCPSocket.Port > 0 {
				lifecycle.PreStop = &corev1.Handler{
					TCPSocket: &corev1.TCPSocketAction{
						Port: intstr.IntOrString{
							IntVal: int32(cc.Lifecycle.PreStop.TCPSocket.Port),
							Type:   intstr.Int,
						},
					},
				}
			}
		}
	}
	return
}
