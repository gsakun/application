package controller

import (
	v3 "github.com/hd-Li/types/apis/project.cattle.io/v3"
	istiorbacv1alpha1 "github.com/hd-Li/types/pkg/istio/apis/rbac/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewServiceRoleObject Use for generate ServiceRoleObject
func NewServiceRoleObject(app *v3.Application) istiorbacv1alpha1.ServiceRole {
	serviceRole := istiorbacv1alpha1.ServiceRole{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(app, v3.SchemeGroupVersion.WithKind("Application"))},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + "servicerole",
			Annotations:     map[string]string{},
		},
		Spec: istiorbacv1alpha1.ServiceRoleSpec{
			Rules: []istiorbacv1alpha1.AccessRule{
				{
					Services: []string{(app.Name + "-" + "service" + "." + app.Namespace + ".svc.cluster.local")},
				},
			},
		},
	}

	return serviceRole
}

// NewServiceRoleBinding Use for generate ServiceRoleBinding
func NewServiceRoleBinding(app *v3.Application) istiorbacv1alpha1.ServiceRoleBinding {
	subjects := []istiorbacv1alpha1.Subject{}

	for _, e := range app.Spec.OptTraits.WhiteList.Users {
		subject := istiorbacv1alpha1.Subject{
			Properties: map[string]string{
				"request.auth.claims[email]": e,
			},
		}

		subjects = append(subjects, subject)
	}

	serviceRoleBinding := istiorbacv1alpha1.ServiceRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceRoleBinding",
			APIVersion: "rbac.istio.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(app, v3.SchemeGroupVersion.WithKind("Application"))},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + "servicerolebinding",
			Annotations:     map[string]string{},
		},
		Spec: istiorbacv1alpha1.ServiceRoleBindingSpec{
			Subjects: subjects,
			RoleRef: istiorbacv1alpha1.RoleRef{
				Kind: "ServiceRole",
				Name: app.Name + "-" + "servicerole",
			},
		},
	}

	return serviceRoleBinding
}
