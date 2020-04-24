package controller

import (
	"fmt"
	"reflect"
	"strings"

	v3 "github.com/hd-Li/types/apis/project.cattle.io/v3"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewAutoScaleInstance(component *v3.Component, app *v3.Application, ref *metav1.OwnerReference) v2beta2.HorizontalPodAutoscaler {
	ownerRef := GetOwnerRef(app)
	split := strings.Split(component.OptTraits.Autoscaling.Metric, "---")
	funcation := string(split[0])
	metric := string(split[1])
	scope := string(split[2])
	value := resource.MustParse(string(component.OptTraits.Autoscaling.Threshold))
	hpa := v2beta2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + component.Version + "hpa",
		},
		Spec: v2beta2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: v2beta2.CrossVersionObjectReference{
				Kind:       "Deployment",
				Name:       app.Name + "-" + component.Name + "-" + "workload" + "-" + component.Version,
				APIVersion: ref.APIVersion,
			},
			MinReplicas: &component.OptTraits.Autoscaling.MinReplicas,
			MaxReplicas: component.OptTraits.Autoscaling.MaxReplicas,
			Metrics: []v2beta2.MetricSpec{
				{
					Type: v2beta2.PodsMetricSourceType,
					Pods: &v2beta2.PodsMetricSource{
						Metric: v2beta2.MetricIdentifier{
							Name: fmt.Sprintf("%s_%s_%s", metric, funcation, scope),
						},
						Target: v2beta2.MetricTarget{
							Type:  v2beta2.AverageValueMetricType,
							Value: &value,
						},
					},
				},
			},
		},
	}
	return hpa
}

func (c *controller) syncHpa(component *v3.Component, app *v3.Application, ref *metav1.OwnerReference) error {
	if !(reflect.DeepEqual(component.OptTraits.Autoscaling, v3.Autoscaling{})) {
		c.syncAutoScaleConfigMap(component, app)
		//c.syncAutoScale(component, app, ref)
	}
	return nil
}

func (c *controller) syncAutoScaleConfigMap(component *v3.Component, app *v3.Application) error {
	log.Infof("Sync autoscaleconfigmap for %s .......\n", app.Namespace+":"+app.Name+"-"+component.Name)
	configmap, err := c.configmapLister.Get("monitor", "adapter-config")
	if err != nil {
		if errors.IsNotFound(err) {
			var stringmap map[string]string = make(map[string]string)
			stringmap["config.yaml"] = ""
			object := NewAutoScaleConfigMapObject(component, app, stringmap)
			_, err = c.configmapClient.Create(&object)
			if err != nil {
				log.Errorf("Create configmap for %s Error : %s\n", "adapter-config", err.Error())
			}
		} else {
			log.Errorf("Get configmap for %s failed", "adapter-config")
		}
	} else {
		var config MetricsDiscoveryConfig
		if configmap != nil {
			value := configmap.Data["config.yml"]
			if value == "" {
				rule := generaterule(component.OptTraits.Autoscaling.Metric, app.Namespace, app.Name+"-"+component.Name+"-"+"workload"+"-"+component.Version)
				config.Rules = append(config.Rules, rule)
			} else {
				data, err := FromYAML([]byte(value))
				if err != nil {
					return err
				} else {
					rule := generaterule(component.OptTraits.Autoscaling.Metric, app.Namespace, app.Name+"-"+component.Name+"-"+"workload"+"-"+component.Version)
					for n, i := range data.Rules {
						if reflect.DeepEqual(i.SeriesQuery, rule.SeriesQuery) {
							break
						} else {
							config.Rules = append(data.Rules[:n], data.Rules[n+1:]...)
							config.Rules = append(data.Rules, rule)
							break
						}
					}
				}
			}
		}
		value, err := yaml.Marshal(config)
		if err != nil {
			return err
		}
		configmap.Data["config.yml"] = string(value)
		_, err = c.configmapClient.Update(configmap)
		if err != nil {
			log.Errorf("Update configmap for %s Error : %s\n", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
			return err
		}
		return nil
	}
	return nil
}

/*func (c *controller) syncAutoScale(component *v3.Component, app *v3.Application, ref *metav1.OwnerReference) error {
	if component.OptTraits.Autoscaling.Metric != "" {
		log.Infof("This app don't need to configure autoscale for %s", app.Namespace+":"+app.Name+"-"+component.Name)
		return nil
	}
	log.Infof("Sync autoscale for %s .......\n", app.Namespace+":"+app.Name+"-"+component.Name)
	insObject := NewAutoScaleInstance(component, app, ref)
	//zk
	insObjectString := GetObjectApplied(insObject)
	insObject.Annotations = make(map[string]string)
	insObject.Annotations[LastAppliedConfigAnnotation] = insObjectString
	instance, err := c.autoscaleLister.Get(app.Namespace, app.Name+"-"+component.Name+"-"+component.Version+"hpa")
	if err != nil {
		log.Debugf("Get autoscale  for %s error : %s\n", (app.Namespace + ":" + app.Name + "-" + component.Name), err.Error())
		if errors.IsNotFound(err) {
			_, err = c.autoscaleClient.Create(&insObject)
			if err != nil {
				log.Errorf("Create autoscale for %s error : %s\n", (app.Namespace + ":" + app.Name + "-" + component.Name), err.Error())
				return nil
			}
		}
	} else {
		if instance != nil {
			if instance.Annotations[LastAppliedConfigAnnotation] != insObjectString {
				insObject.ObjectMeta.ResourceVersion = instance.ObjectMeta.ResourceVersion
				_, err = c.autoscaleClient.Update(&insObject)
				if err != nil {
					log.Errorf("Update autoscale for %s error : %s\n", (app.Namespace + ":" + app.Name + "-" + component.Name), err.Error())
				}
			}
		}
	}
	return nil
}*/

func NewAutoScaleConfigMapObject(component *v3.Component, app *v3.Application, data map[string]string) corev1.ConfigMap {
	configmap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "adapter-config",
			Name:      "monitor",
		},
		Data: data,
	}
	return configmap
}

func FromYAML(contents []byte) (*MetricsDiscoveryConfig, error) {
	var cfg MetricsDiscoveryConfig
	if err := yaml.UnmarshalStrict(contents, &cfg); err != nil {
		return nil, fmt.Errorf("unable to parse metrics discovery config: %v", err)
	}
	return &cfg, nil
}

func generaterule(data, namespace, podnameprefix string) (rule DiscoveryRule) {
	split := strings.Split(data, "---")
	funcation := string(split[0])
	metric := string(split[1])
	scope := string(split[2])
	time := string(split[3])
	var rmap map[string]GroupResource = make(map[string]GroupResource)
	rmap["kubernetes_namespace"] = GroupResource{
		Resource: "namespace",
	}
	rmap["kubernetes_pod_name"] = GroupResource{
		Resource: "pod",
	}
	//rule.SeriesQuery = fmt.Sprintf(`%s{kubernetes_namespace="%s",kubernetes_pod_name=~"%s.*"}`, app.Namespace, app.Name+"-"+component.Name+"-"+"workload"+"-"+component.Version, metric)
	rule.SeriesQuery = fmt.Sprintf(`%s{kubernetes_namespace="%s",kubernetes_pod_name=~"%s.*"}`, metric, namespace, podnameprefix)
	rule.Resources = ResourceMapping{
		Overrides: rmap,
	}
	rule.Name = NameMapping{
		Matches: metric,
		As:      fmt.Sprintf("${1}_%s_%s", funcation, scope),
	}
	if scope == "all" {
		rule.MetricsQuery = fmt.Sprintf("%s(<<.Series>>{<<.LabelMatchers>>}[%s])", metric, time)
	}
	if scope == "per" {
		rule.MetricsQuery = fmt.Sprintf("%s(<<.Series>>{<<.LabelMatchers>>}[%s]) by (<<.GroupBy>>)", metric, time)
	}
	return
}
