package controller

import (
	v3 "github.com/hd-Li/types/apis/project.cattle.io/v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetConfigMap(conf v3.ConfigFile, container *v3.ComponentContainer, component *v3.Component, app *v3.Application) corev1.ConfigMap {
	ownerRef := GetOwnerRef(app)

	configMap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + container.Name + "-" + "conf",
		},
		Data: map[string]string{
			conf.FileName: conf.Value,
		},
	}

	return configMap
}
