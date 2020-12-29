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
	if app.Spec.OptTraits.WhiteList != nil {
		for _, e := range RemoveRepByLoop(app.Spec.OptTraits.WhiteList.Users) {
			subject := istiorbacv1alpha1.Subject{
				Properties: map[string]string{
					"request.auth.claims[email]": e,
				},
			}

			subjects = append(subjects, subject)
		}
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
			RoleRef: istiorbacv1alpha1.RoleRef{
				Kind: "ServiceRole",
				Name: app.Name + "-" + "servicerole",
			},
		},
	}
	if len(subjects) != 0 {
		serviceRoleBinding.Spec.Subjects = subjects
	}

	return serviceRoleBinding
}

// RemoveRepByLoop use for remove Rep element
func RemoveRepByLoop(slc []string) []string {
	result := []string{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}
