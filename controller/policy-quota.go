package controller

import (
	"os"

	v3 "github.com/hd-Li/types/apis/project.cattle.io/v3"
	"github.com/hd-Li/types/pkg/istio/apis/config/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewQuotaInstance Use for generate QuotaInstance
func NewQuotaInstance(app *v3.Application) v1alpha2.Instance {
	instance := v1alpha2.Instance{
		TypeMeta: metav1.TypeMeta{
			Kind:       "instance",
			APIVersion: "config.istio.io/v1alpha2",
		},
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(app, v3.SchemeGroupVersion.WithKind("Application"))},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + "quotainstance",
			Annotations:     map[string]string{},
		},
		Spec: v1alpha2.InstanceSpec{
			CompiledTemplate: "quota",
			Params: v1alpha2.InstanceParams{
				Dimensions: map[string]string{
					"destination":               `destination.labels["app"] | destination.workload.name | "unknown"`,
					"destinationVersion":        `destination.labels["version"] | "unknown"`,
					"request_auth_claims_email": `request.auth.claims["email"] | "unknown"`,
					"source":                    `request.headers["x-forwarded-for"] | "unknown"`,
				},
			},
		},
	}

	return instance
}

// NewQuotaSpec Use for generate QuotaSpec
func NewQuotaSpec(app *v3.Application) v1alpha2.QuotaSpec {
	quota := v1alpha2.Quota{
		Quota:  app.Name + "-" + "quotainstance",
		Charge: 1,
	}

	quotaRule := v1alpha2.QuotaRule{
		Quotas: []*v1alpha2.Quota{&quota},
	}

	quotaspec := v1alpha2.QuotaSpec{
		TypeMeta: metav1.TypeMeta{
			Kind:       "QuotaSpec",
			APIVersion: "config.istio.io/v1alpha2",
		},
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(app, v3.SchemeGroupVersion.WithKind("Application"))},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + "quotaspec",
			Annotations:     map[string]string{},
		},
		Spec: v1alpha2.QuotaSubSpec{
			Rules: []*v1alpha2.QuotaRule{&quotaRule},
		},
	}

	return quotaspec
}

// NewQuotaSpecBinding Use for generate QuotaSpecBinding
func NewQuotaSpecBinding(app *v3.Application) v1alpha2.QuotaSpecBinding {
	istioService := v1alpha2.IstioService{
		Name:      app.Name + "-" + "service",
		Namespace: app.Namespace,
	}

	quotaSpecReference := v1alpha2.QuotaSpecBindingQuotaSpecReference{
		Name:      app.Name + "-" + "quotaspec",
		Namespace: app.Namespace,
	}

	quotaspecbinding := v1alpha2.QuotaSpecBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "QuotaSpecBinding",
			APIVersion: "config.istio.io/v1alpha2",
		},
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(app, v3.SchemeGroupVersion.WithKind("Application"))},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + "quotaspecbinding",
			Annotations:     map[string]string{},
		},
		Spec: v1alpha2.QuotaSpecBindingSpec{
			Services:   []*v1alpha2.IstioService{&istioService},
			QuotaSpecs: []*v1alpha2.QuotaSpecBindingQuotaSpecReference{&quotaSpecReference},
		},
	}

	return quotaspecbinding
}

// NewQuotaHandlerObject Use for generate QuotaHandlerObject
func NewQuotaHandlerObject(app *v3.Application) *v1alpha2.Handler {
	redisServer := os.Getenv("REDIS_SERVER")

	overrides := []v1alpha2.Override{}
	for _, v := range app.Spec.OptTraits.RateLimit.Overrides {
		override := v1alpha2.Override{
			MaxAmount: v.RequestAmount,
			Dimensions: map[string]string{
				"request_auth_claims_email": v.User,
			},
		}
		overrides = append(overrides, override)
	}

	handlerquota := v1alpha2.HandlerQuota{
		Name:               app.Name + "-" + "quotainstance" + "." + "instance" + "." + app.Namespace,
		MaxAmount:          app.Spec.OptTraits.RateLimit.RequestAmount,
		ValidDuration:      app.Spec.OptTraits.RateLimit.TimeDuration,
		BucketDuration:     "200ms",
		RateLimitAlgorithm: v1alpha2.ROLLING,
		Overrides:          overrides,
	}

	handler := v1alpha2.Handler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "handler",
			APIVersion: "config.istio.io/v1alpha2",
		},
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(app, v3.SchemeGroupVersion.WithKind("Application"))},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + "quotahandler",
			Annotations:     map[string]string{},
		},
		Spec: v1alpha2.HandlerSpec{
			CompiledAdapter: "redisquota",
			Params: v1alpha2.HandlerParams{
				RedisServerUrl:     redisServer,
				ConnectionPoolSize: 10,
				Quotas:             []v1alpha2.HandlerQuota{handlerquota},
			},
		},
	}

	return &handler
}

// NewQuotaRuleObject Use for generate QuotaRuleObject
func NewQuotaRuleObject(app *v3.Application) v1alpha2.Rule {
	instance := app.Name + "-" + "quotainstance" + "." + "instance" + "." + app.Namespace
	action := v1alpha2.Action{
		Handler:   app.Name + "-" + "quotahandler" + "." + "handler" + "." + app.Namespace,
		Instances: []string{instance},
	}

	rule := v1alpha2.Rule{
		TypeMeta: metav1.TypeMeta{
			Kind:       "rule",
			APIVersion: "config.istio.io/v1alpha2",
		},
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(app, v3.SchemeGroupVersion.WithKind("Application"))},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + "quotarule",
			Annotations:     map[string]string{},
		},
		Spec: v1alpha2.RuleSpec{
			Actions: []*v1alpha2.Action{&action},
		},
	}

	return rule
}
