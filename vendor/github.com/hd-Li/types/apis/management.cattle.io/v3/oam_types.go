package v3

import (
	projectv3 "github.com/hd-Li/types/apis/project.cattle.io/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
apiVersion: core.oam.dev/v1alpha1
kind: ComponentSchematic
metadata:
  name: frontend
  annotations:
    version: "1.0.0"
    description: Node.js Web Server
spec:
  workloadType: core.oam.dev/v1alpha1.Server
  parameters:
    - name: database
      type: string
      required: false
  containers:
    - name: frontend
      ports:
        - containerPort: 3000
          name: http
      image: janakiramm/todo:v1
      env:
        - name: DB
          value: db
          fromParam: database*/

type AppTraits struct{}

type ApplicationConfigurationSpec struct {
	Parameters []projectv3.Parameter `json:"parameters,omitempty"`
	Components []projectv3.Component `json:"components"`
	AppTraits  AppTraits             `json:"appTraits,omitempty"`
}

type ApplicationConfigurationTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationConfigurationSpec   `json:"spec"`
	Status ApplicationConfigurationStatus `json:"status,omitempty"`
}

type ApplicationConfigurationStatus struct{}
