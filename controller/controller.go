package controller

import (
	"context"
	"strings"

	"reflect"

	log "github.com/sirupsen/logrus"

	//typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"github.com/hd-Li/types/apis/apps/v1beta2"
	"github.com/hd-Li/types/apis/autoscaling/v2beta2"
	v1 "github.com/hd-Li/types/apis/core/v1"
	"github.com/hd-Li/types/config"
	"k8s.io/apimachinery/pkg/runtime"

	//utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	v3 "github.com/hd-Li/types/apis/project.cattle.io/v3"
	//appsv1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/apimachinery/pkg/api/errors"
	//"k8s.io/client-go/kubernetes/scheme"
	istioauthnv1alpha1 "github.com/hd-Li/types/apis/authentication.istio.io/v1alpha1"
	istioconfigv1alpha2 "github.com/hd-Li/types/apis/config.istio.io/v1alpha2"
	istionetworkingv1alph3 "github.com/hd-Li/types/apis/networking.istio.io/v1alpha3"
	istiorbacv1alpha1 "github.com/hd-Li/types/apis/rbac.istio.io/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	//istiov1alpha3 "github.com/knative/pkg/apis/istio/v1alpha3"
)

const (
	// LastAppliedConfigAnnotation define LastAppliedConfigAnnotation
	LastAppliedConfigAnnotation string = "application/last-applied-configuration"
)

type controller struct {
	applicationClient        v3.ApplicationInterface
	applicationLister        v3.ApplicationLister
	nsClient                 v1.NamespaceInterface
	coreV1                   v1.Interface
	appsV1beta2              v1beta2.Interface
	podLister                v1.PodLister                             //zk
	podClient                v1.PodInterface                          //zk
	secretLister             v1.SecretLister                          //zk
	secretClient             v1.SecretInterface                       //zk
	autoscaleLister          v2beta2.HorizontalPodAutoscalerLister    //zk
	autoscaleClient          v2beta2.HorizontalPodAutoscalerInterface //zk
	deploymentLister         v1beta2.DeploymentLister
	deploymentClient         v1beta2.DeploymentInterface
	serviceLister            v1.ServiceLister
	serviceClient            v1.ServiceInterface
	virtualServiceLister     istionetworkingv1alph3.VirtualServiceLister
	virtualServiceClient     istionetworkingv1alph3.VirtualServiceInterface
	destLister               istionetworkingv1alph3.DestinationRuleLister
	destClient               istionetworkingv1alph3.DestinationRuleInterface
	configmapLister          v1.ConfigMapLister    //zk
	configmapClient          v1.ConfigMapInterface //zk
	gatewayLister            istionetworkingv1alph3.GatewayLister
	gatewayClient            istionetworkingv1alph3.GatewayInterface
	policyLister             istioauthnv1alpha1.PolicyLister
	policyClient             istioauthnv1alpha1.PolicyInterface
	clusterconfigLister      istiorbacv1alpha1.ClusterRbacConfigLister
	clusterconfigClient      istiorbacv1alpha1.ClusterRbacConfigInterface
	serviceRoleLister        istiorbacv1alpha1.ServiceRoleLister
	serviceRoleClient        istiorbacv1alpha1.ServiceRoleInterface
	serviceRoleBindingLister istiorbacv1alpha1.ServiceRoleBindingLister
	serviceRoleBindingClient istiorbacv1alpha1.ServiceRoleBindingInterface
	handerLister             istioconfigv1alpha2.HandlerLister
	handlerClient            istioconfigv1alpha2.HandlerInterface
	ruleLister               istioconfigv1alpha2.RuleLister
	ruleClient               istioconfigv1alpha2.RuleInterface
	instanceLister           istioconfigv1alpha2.InstanceLister
	instanceClient           istioconfigv1alpha2.InstanceInterface
	quotaspecLister          istioconfigv1alpha2.QuotaSpecLister
	quotaspecClient          istioconfigv1alpha2.QuotaSpecInterface
	quotaspecbindingLister   istioconfigv1alpha2.QuotaSpecBindingLister
	quotaspecbindingClient   istioconfigv1alpha2.QuotaSpecBindingInterface
	recorder                 record.EventRecorder
}

// Register all resource
func Register(ctx context.Context, userContext *config.UserOnlyContext) {
	/*
		utilruntime.Must(v3.AddToScheme(scheme.Scheme))
		log.Infoln("Creating event broadcaster")
		eventBroadcaster := record.NewBroadcaster()
		//eventBroadcaster.StartLogging(fmt.Printf)
		eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: userContext.Core.Events("")})
		recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "application-controllere"})
	*/
	c := controller{
		applicationClient:        userContext.Project.Applications(""),
		applicationLister:        userContext.Project.Applications("").Controller().Lister(),
		nsClient:                 userContext.Core.Namespaces(""),
		coreV1:                   userContext.Core,
		appsV1beta2:              userContext.Apps,
		deploymentLister:         userContext.Apps.Deployments("").Controller().Lister(),
		deploymentClient:         userContext.Apps.Deployments(""),
		configmapLister:          userContext.Core.ConfigMaps("").Controller().Lister(), //zk
		configmapClient:          userContext.Core.ConfigMaps(""),                       //zk
		podLister:                userContext.Core.Pods("").Controller().Lister(),       //zk
		podClient:                userContext.Core.Pods(""),                             //zk
		secretLister:             userContext.Core.Secrets("").Controller().Lister(),    //zk
		secretClient:             userContext.Core.Secrets(""),                          //zk
		serviceLister:            userContext.Core.Services("").Controller().Lister(),
		serviceClient:            userContext.Core.Services(""),
		autoscaleLister:          userContext.Autoscaling.HorizontalPodAutoscalers("").Controller().Lister(), //zk
		autoscaleClient:          userContext.Autoscaling.HorizontalPodAutoscalers(""),                       //zk
		virtualServiceLister:     userContext.IstioNetworking.VirtualServices("").Controller().Lister(),
		virtualServiceClient:     userContext.IstioNetworking.VirtualServices(""),
		destLister:               userContext.IstioNetworking.DestinationRules("").Controller().Lister(),
		destClient:               userContext.IstioNetworking.DestinationRules(""),
		gatewayLister:            userContext.IstioNetworking.Gateways("").Controller().Lister(),
		gatewayClient:            userContext.IstioNetworking.Gateways(""),
		policyLister:             userContext.IstioAuthn.Policies("").Controller().Lister(),
		policyClient:             userContext.IstioAuthn.Policies(""),
		clusterconfigLister:      userContext.IstioRbac.ClusterRbacConfigs("").Controller().Lister(),
		clusterconfigClient:      userContext.IstioRbac.ClusterRbacConfigs(""),
		serviceRoleLister:        userContext.IstioRbac.ServiceRoles("").Controller().Lister(),
		serviceRoleClient:        userContext.IstioRbac.ServiceRoles(""),
		serviceRoleBindingLister: userContext.IstioRbac.ServiceRoleBindings("").Controller().Lister(),
		serviceRoleBindingClient: userContext.IstioRbac.ServiceRoleBindings(""),
		handerLister:             userContext.IstioConfig.Handlers("").Controller().Lister(),
		handlerClient:            userContext.IstioConfig.Handlers(""),
		ruleLister:               userContext.IstioConfig.Rules("").Controller().Lister(),
		ruleClient:               userContext.IstioConfig.Rules(""),
		instanceLister:           userContext.IstioConfig.Instances("").Controller().Lister(),
		instanceClient:           userContext.IstioConfig.Instances(""),
		quotaspecLister:          userContext.IstioConfig.QuotaSpecs("").Controller().Lister(),
		quotaspecClient:          userContext.IstioConfig.QuotaSpecs(""),
		quotaspecbindingLister:   userContext.IstioConfig.QuotaSpecBindings("").Controller().Lister(),
		quotaspecbindingClient:   userContext.IstioConfig.QuotaSpecBindings(""),
	}

	c.applicationClient.AddHandler(ctx, "applictionCreateOrUpdate", c.sync)
}

func (c *controller) sync(key string, app *v3.Application) (runtime.Object, error) {
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
	if app == nil {
		return nil, nil
	}
	//log.Infof("application info %v", application)
	//app := application.DeepCopy()

	c.syncNamespaceCommon(app)

	//the deployed app is trusted or not
	var trusted bool = false

	components := app.Spec.Components
	if len(components) == 0 && len(app.Status.ComponentResource) == 0 {
		return nil, nil
	}
	//zk
	var oldcomresource map[string]v3.ComponentResources = make(map[string]v3.ComponentResources)
	if app.Status.ComponentResource != nil {
		oldcomresource = app.DeepCopy().Status.ComponentResource
	}
	app.Status.ComponentResource = make(map[string]v3.ComponentResources)
	/*if reflect.DeepEqual(app.Spec.OptTraits.ImagePullConfig, v3.ImagePullConfig{}) {
		log.Debugf("Component %s's image pull configuration is not configured,ignore it")
	} else {
		if app.Spec.OptTraits.ImagePullConfig.Username == "" || app.Spec.OptTraits.ImagePullConfig.Registry == "" || app.Spec.OptTraits.ImagePullConfig.Password == "" {
			log.Errorf("Component %s's imagepullconfig need review", app.Name)
		} else {
			_, _ = c.syncImagePullSecrets(app)
		}
	}*/ //Not needed for the time being
	var deletelist []string
	for _, component := range components {
		//if containers is nil, the app is trusted, this controller does not manage its workload's lifecycle
		if len(component.Containers) == 0 {
			trusted = true
		}
		ownerRefOfDeploy := new(metav1.OwnerReference)
		if trusted == false {
			delete(oldcomresource, (app.Name + "_" + component.Name + "_" + component.Version))
			c.syncConfigmaps(&component, app)
			err := c.syncWorkload(&component, app, ownerRefOfDeploy)
			if err != nil {
				return nil, err
			}
		} else {
			err := c.syncTrustedWorkload(&component, app, ownerRefOfDeploy)
			if err != nil {
				return nil, err
			}
		}
		//log.Infof("ownerRefOfDeploy INFO IS %v", ownerRefOfDeploy)
		if ownerRefOfDeploy.APIVersion != "" {
			c.syncHpa(&component, app, ownerRefOfDeploy)
		}
	}
	if len(app.Spec.OptTraits.Fusing.PodList) != 0 {
		log.Infoln("START FUSING")
		var action bool = false
		if app.Spec.OptTraits.Fusing.Action == "in" {
			action = true
		}
		for _, i := range app.Spec.OptTraits.Fusing.PodList {
			c.syncFusing(i, app.Namespace, action)
		}
		app.Spec.OptTraits.Fusing = v3.Fusing{}
	}
	c.syncService(app)
	c.syncAuthor(app)
	c.syncPolicy(app)
	log.Debugf("These versions need to be removed %v", oldcomresource)
	for k := range oldcomresource {
		deletelist = append(deletelist, k)
	}
	if len(deletelist) != 0 {
		errlist := c.gc(app.Namespace, deletelist)
		if len(errlist) != 0 {
			for _, i := range errlist {
				app.Status.ComponentResource[i] = v3.ComponentResources{}
			}
		}
	}
	c.syncStatus(app)
	return nil, nil
}

func (c *controller) syncNamespaceCommon(app *v3.Application) error {
	log.Infof("Sync namespaceCommon for %s", app.Namespace+":"+app.Name)

	var ns *corev1.Namespace
	var err error

	for i := 0; i < 3; i++ {
		ns, err = c.nsClient.Get(app.Namespace, metav1.GetOptions{})
		if err != nil {
			log.Infof("Get namespace object error for app %s error : %s", (app.Namespace + ":" + app.Name), err.Error())
		} else {
			break
		}
	}
	_, err = c.gatewayLister.Get(app.Namespace, (app.Namespace + "-" + "gateway"))
	if err != nil {
		log.Infof("Get gateway error for %s error : %s", (app.Namespace + ":" + app.Name), err.Error())

		if errors.IsNotFound(err) {
			gateway := NewGatewayObject(app, ns)
			_, err = c.gatewayClient.Create(&gateway)
			if err != nil {
				log.Infof("Create gateway error for %s error : %s", (app.Namespace + ":" + app.Name), err.Error())
			}
		}
	}
	log.Infof("Sync gateway done for namespace %s", app.Namespace)

	_, err = c.policyLister.Get(app.Namespace, "default")
	if err != nil {
		log.Infof("Get policy for %s error : %s", (app.Namespace + ":" + app.Name), err.Error())
		if errors.IsNotFound(err) {
			policy := NewPolicyObject(app, ns)
			_, err = c.policyClient.Create(&policy)
			if err != nil {
				log.Infof("Create policy error for %s error : %s", (app.Namespace + ":" + app.Name), err.Error())
			}
		}
	}
	log.Infof("Sync policy done for %s", app.Namespace)

	cfg, err := c.clusterconfigLister.Get("", "default")
	if err != nil {
		log.Infof("Get clusterrbacconfig for %s error : %s", (app.Namespace + ":" + app.Name), err.Error())
		if errors.IsNotFound(err) {
			clusterConfig := NewClusterRbacConfig(app, ns)
			_, err = c.clusterconfigClient.Create(&clusterConfig)
			if err != nil {
				log.Errorf("Create clusterrbacconfig error for %s error : %s", (app.Namespace + ":" + app.Name), err.Error())
			}
		}
	} else {
		if cfg != nil {
			clusterrbacconfig := cfg.DeepCopy()
			if _, ok := clusterrbacconfig.ObjectMeta.Labels[app.Namespace]; !ok {
				clusterrbacconfig.Spec.Inclusion.Namespaces = append(clusterrbacconfig.Spec.Inclusion.Namespaces, app.Namespace)
				clusterrbacconfig.ObjectMeta.Labels[app.Namespace] = "included"
				clusterrbacconfig.Namespace = "default" //avoid the client-go bug
				_, err = c.clusterconfigClient.Update(clusterrbacconfig)
				if err != nil {
					log.Errorf("Update clusterrbacconfig error for %s error : %s", (app.Namespace + ":" + app.Name), err.Error())
				}
			}
		}
	}
	log.Infof("Sync clusterrbacconfig done for %s", app.Namespace)

	return nil
}

func (c *controller) syncConfigmaps(component *v3.Component, app *v3.Application) error {
	log.Infof("Sync configmap for %s", app.Namespace+":"+component.Name+":"+component.Version)
	object := NewConfigMapObject(component, app)
	if len(object.Data) == 0 {
		log.Debugf("ConfigMap data is nil, Do not need sync configmap for %s", app.Namespace+":"+app.Name+":"+component.Name+":"+component.Version)
		return nil
	}
	appliedString := GetObjectApplied(object)
	configmapname := app.Name + "-" + component.Name + "-" + component.Version + "-" + "configmap"
	configmap, err := c.configmapLister.Get(app.Namespace, configmapname)
	object.Annotations = make(map[string]string)
	object.Annotations[LastAppliedConfigAnnotation] = appliedString
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = c.configmapClient.Create(&object)
			if err != nil {
				log.Errorf("Create configmap for %s Error : %s", (app.Namespace + ":" + app.Name + ":" + component.Name + ":" + component.Version), err.Error())
			}
		} else {
			log.Errorf("Get configmap for %s failed", configmapname)
		}
	} else {
		if configmap != nil {
			if configmap.Annotations[LastAppliedConfigAnnotation] != appliedString {
				_, err := c.configmapClient.Update(&object)
				if err != nil {
					log.Errorf("Update configmap for %s Error : %s", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
					return nil
				}
			}
		}
	}
	return nil
}

func (c *controller) syncImagePullSecrets(app *v3.Application) (string, error) {
	log.Infof("Sync imagepull secret for %s", app.Namespace)
	object := NewSecretObject(app)
	if len(object.Data) == 0 {
		log.Debugf("ImagePullSecret not define,ignore %s", app.Namespace+":"+app.Name)
		return "", nil
	}
	appliedString := GetObjectApplied(object)

	secretname := app.Name + "-" + "registry-secret"
	secret, err := c.secretLister.Get(app.Namespace, secretname)
	object.Annotations = make(map[string]string)
	object.Annotations[LastAppliedConfigAnnotation] = appliedString
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = c.secretClient.Create(&object)
			if err != nil {
				log.Errorf("Create secret for %s Error : %s", (app.Namespace + ":" + app.Name), err.Error())
				return "", err
			}
			log.Infof("Create Secret %s successful", secretname)
			return secretname, nil

		}
		log.Errorf("Get sercret for %s failed", secretname)
		return "", err
	}
	if secret != nil {
		if secret.Annotations[LastAppliedConfigAnnotation] != appliedString {
			_, err := c.secretClient.Update(&object)
			if err != nil {
				log.Errorf("Update secret for %s Error : %s", (app.Namespace + ":" + app.Name), err.Error())
				return "", err
			}
			return secretname, nil
		}
	}
	return secretname, nil
}

func (c *controller) syncWorkload(component *v3.Component, app *v3.Application, ref *metav1.OwnerReference) error {
	resourceWorkloadType := "deployment"
	if resourceWorkloadType == "deployment" {
		err := c.syncDeployment(component, app, ref)
		if err != nil {
			return err
		}
	}

	return nil
}

//zk
func (c *controller) syncStatus(app *v3.Application) {
	_, err := c.applicationClient.Update(app)
	if err != nil {
		log.Errorf("Update application for %s Error : %s", (app.Namespace + ":" + app.Name), err.Error())
	} else {
		log.Infof("Update application for %s Done", (app.Namespace + ":" + app.Name))
	}
}

func (c *controller) syncDeployment(component *v3.Component, app *v3.Application, ref *metav1.OwnerReference) error {
	log.Infof("Sync deploy for %s", app.Namespace+":"+component.Name)
	object := NewDeployObject(component, app)
	appliedString := GetObjectApplied(object)
	//zk
	object.Annotations = make(map[string]string)
	object.Annotations[LastAppliedConfigAnnotation] = appliedString
	deploy, err := c.deploymentLister.Get(app.Namespace, app.Name+"-"+component.Name+"-"+"workload"+"-"+component.Version)
	if err != nil {
		//log.Infof("Get deploy for %s Error : %s", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
		if errors.IsNotFound(err) {
			getdeploy, err := c.deploymentClient.Create(&object)
			if err != nil {
				log.Errorf("Create deploy for %s Error : %s", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
				return err
			}
			*ref = *(metav1.NewControllerRef(getdeploy, v1beta2.SchemeGroupVersion.WithKind("Deployment")))
		}
	} else {
		if deploy != nil {
			if deploy.Annotations[LastAppliedConfigAnnotation] != appliedString {
				getdeploy, err := c.deploymentClient.Update(&object)
				if err != nil {
					log.Errorf("Update deploy for %s Error : %s", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
					return err
				}
				*ref = *(metav1.NewControllerRef(getdeploy, v1beta2.SchemeGroupVersion.WithKind("Deployment")))
			} else {
				*ref = *(metav1.NewControllerRef(deploy, v1beta2.SchemeGroupVersion.WithKind("Deployment")))
			}
		}
	}
	log.Infof("Sync deploy for %s done!", app.Namespace+":"+app.Name+":"+component.Name)
	app.Status.ComponentResource[(app.Name + "_" + component.Name + "_" + component.Version)] = v3.ComponentResources{
		Workload: object.Name,
	}
	return nil
}

func (c *controller) syncService(app *v3.Application) error {
	log.Infof("Sync service for %s", app.Name)
	object := NewServiceObject(app)
	object.ObjectMeta.OwnerReferences = []metav1.OwnerReference{*metav1.NewControllerRef(app, v3.SchemeGroupVersion.WithKind("Application"))}
	objectString := GetObjectApplied(object)
	//zk
	object.Annotations = make(map[string]string)
	object.Annotations[LastAppliedConfigAnnotation] = objectString

	service, err := c.serviceLister.Get(app.Namespace, app.Name+"-"+"service")
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = c.serviceClient.Create(&object)
			if err != nil {
				log.Errorf("Create service for %s Error : %s", (app.Namespace + ":" + app.Name), err.Error())
			}
		}
	} else {
		if service != nil {
			if service.Annotations[LastAppliedConfigAnnotation] != objectString {
				//c.serviceClient.DeleteNamespaced(service.Namespace, service.Name, &metav1.DeleteOptions{})
				_, err = c.serviceClient.Update(&object)
				if err != nil {
					log.Errorf("Update(Create) Service for %s Error : %s", (app.Namespace + ":" + app.Name), err.Error())
				}
			}
		}

		_, err = c.serviceRoleLister.Get(app.Namespace, app.Name+"-"+"servicerole")
		if err != nil {
			if errors.IsNotFound(err) {
				svcRoleObject := NewServiceRoleObject(app)
				_, err = c.serviceRoleClient.Create(&svcRoleObject)
				if err != nil {
					log.Errorf("Create ServiceRole for %s Error : %s", (app.Name), err.Error())
				}
			}
		}
	}
	vsObject := NewVirtualServiceObject(app)
	vsObjectString := GetObjectApplied(vsObject)
	vsObject.Annotations[LastAppliedConfigAnnotation] = vsObjectString

	vs, err := c.virtualServiceLister.Get(app.Namespace, (app.Name + "-" + "vs"))
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = c.virtualServiceClient.Create(&vsObject)
			if err != nil {
				log.Errorf("Create VirtualService error for %s error : %s", (app.Namespace + ":" + app.Name), err.Error())
			}
		}
	} else {
		if vs != nil {
			if vs.Annotations[LastAppliedConfigAnnotation] != vsObjectString {
				vsObject.ObjectMeta.ResourceVersion = vs.ObjectMeta.ResourceVersion
				_, err = c.virtualServiceClient.Update(&vsObject)
				if err != nil {
					log.Errorf("Update VirtualService error for %s error : %s", (app.Namespace + ":" + app.Name), err.Error())
				}
			}
		}
	}
	if !(reflect.DeepEqual(app.Spec.OptTraits.LoadBalancer, v3.LoadBalancerSettings{})) || !(reflect.DeepEqual(app.Spec.OptTraits.CircuitBreaking, v3.CircuitBreaking{})) {
		destObject := NewDestinationruleObject(app)
		destObjectString := GetObjectApplied(destObject)
		destObject.Annotations[LastAppliedConfigAnnotation] = destObjectString

		dest, err := c.destLister.Get(app.Namespace, (app.Name + "-" + "destinationrule"))
		if err != nil {
			//log.Errorf("Get DestinationRule error for %s error : %s", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
			if errors.IsNotFound(err) {
				_, err = c.destClient.Create(&destObject)
				if err != nil {
					log.Errorf("Create DestinationRule error for %s error : %s", (app.Namespace + ":" + app.Name), err.Error())
				}
			}
		} else {

			if dest != nil {
				if dest.Annotations[LastAppliedConfigAnnotation] != destObjectString {
					destObject.ObjectMeta.ResourceVersion = dest.ObjectMeta.ResourceVersion
					_, err := c.destClient.Update(&destObject)
					if err != nil {
						log.Errorf("Update DestinationRule error for %s error : %s", (app.Namespace + ":" + app.Name), err.Error())
					}
				}
			}
		}
	}

	return nil
}
func (c *controller) syncAuthor(app *v3.Application) error {
	object := NewServiceRoleBinding(app)
	objectString := GetObjectApplied(object)
	object.Annotations = make(map[string]string)
	object.Annotations[LastAppliedConfigAnnotation] = objectString

	serviceRoleBinding, err := c.serviceRoleBindingLister.Get(app.Namespace, object.Name)
	if err != nil {
		if errors.IsNotFound(err) {
			if len(app.Spec.OptTraits.WhiteList.Users) == 0 {
				log.Infoln("whitelist.user is nil,there is nothing to do")
				return nil
			}
			_, err = c.serviceRoleBindingClient.Create(&object)
			if err != nil {
				log.Errorf("Create servicerolebinding error for %s error : %s", (app.Namespace + ":" + app.Name), err.Error())
			}
		}
	} else {
		if serviceRoleBinding != nil {
			if serviceRoleBinding.Annotations[LastAppliedConfigAnnotation] != objectString {
				if len(app.Spec.OptTraits.WhiteList.Users) == 0 {
					log.Infof("whitelist is null ,need delete servicerolebinding and servicerole for %s", app.Name)
					err = c.serviceRoleBindingClient.DeleteNamespaced(app.Namespace, app.Name+"-"+"servicerolebinding", &metav1.DeleteOptions{})
					if err != nil {
						log.Errorln(err)
					}
					err = c.serviceRoleClient.DeleteNamespaced(app.Namespace, app.Name+"-"+"servicerole", &metav1.DeleteOptions{})
					if err != nil {
						log.Errorln(err)
					}
					return nil
				}
				//object.ObjectMeta.ResourceVersion = serviceRoleBinding.ObjectMeta.ResourceVersion
				_, err = c.serviceRoleBindingClient.Update(&object)
				if err != nil {
					log.Errorf("Update servicerolebinding error for %s error : %s", (app.Namespace + ":" + app.Name), err.Error())
				}
			}
		}
	}
	return nil
}

func (c *controller) syncPolicy(app *v3.Application) error {
	if app.Spec.OptTraits.RateLimit.TimeDuration != "" {
		c.syncQuotaPolicy(app)
	}
	return nil
}

func (c *controller) syncQuotaPolicy(app *v3.Application) error {
	log.Infof("Sync quotapolicy for %s", app.Namespace+":"+app.Name)

	insObject := NewQuotaInstance(app)
	//zk
	insObjectString := GetObjectApplied(insObject)
	insObject.Annotations = make(map[string]string)
	insObject.Annotations[LastAppliedConfigAnnotation] = insObjectString

	instance, err := c.instanceLister.Get(app.Namespace, app.Name+"-"+"quotainstance")
	if err != nil {
		//log.Infof("Get quotapolicy  for %s error : %s", (app.Namespace + ":" + app.Name + "-" + component.Name), err.Error())
		if errors.IsNotFound(err) {
			_, err = c.instanceClient.Create(&insObject)
			if err != nil {
				log.Errorf("Create quotapolicy  for %s error : %s", (app.Namespace + ":" + app.Name), err.Error())
				return nil
			}
		}
	} else {

		if instance != nil {
			if instance.Annotations[LastAppliedConfigAnnotation] != insObjectString {
				insObject.ObjectMeta.ResourceVersion = instance.ObjectMeta.ResourceVersion
				_, err = c.instanceClient.Update(&insObject)
				if err != nil {
					log.Errorf("Update quotapolicy  for %s error : %s", (app.Namespace + ":" + app.Name), err.Error())
				}
			}
		}
	}
	//config for client
	specObject := NewQuotaSpec(app)
	specObjectString := GetObjectApplied(specObject)
	specObject.Annotations = make(map[string]string)
	specObject.Annotations[LastAppliedConfigAnnotation] = specObjectString

	_, err = c.quotaspecLister.Get(app.Namespace, app.Name+"-"+"quotaspec")
	if err != nil {
		//log.Infof("Get quotaspec  for %s error : %s", (app.Namespace + ":" + app.Name + "-" + component.Name), err.Error())
		if errors.IsNotFound(err) {
			_, err = c.quotaspecClient.Create(&specObject)
			if err != nil {
				log.Errorf("Create quotaspec  for %s error : %s", (app.Namespace + ":" + app.Name), err.Error())
				return nil
			}
		}
	}

	specbindingObject := NewQuotaSpecBinding(app)
	specbindingObjectString := GetObjectApplied(specbindingObject)
	specbindingObject.Annotations = make(map[string]string)
	specbindingObject.Annotations[LastAppliedConfigAnnotation] = specbindingObjectString

	_, err = c.quotaspecbindingLister.Get(app.Namespace, app.Name+"-"+"quotaspecbinding")
	if err != nil {
		//log.Errorf("Get quotaspecbinding for %s error : %s", (app.Namespace + ":" + app.Name + "-" + component.Name), err.Error())
		if errors.IsNotFound(err) {
			_, err = c.quotaspecbindingClient.Create(&specbindingObject)
			if err != nil {
				log.Errorf("Create quotaspecbinding  for %s error : %s", (app.Namespace + ":" + app.Name), err.Error())
				return nil
			}
		}
	}

	//config for (mixer) server
	qhObject := NewQuotaHandlerObject(app)
	qhObjectString := GetObjectApplied(qhObject)
	qhObject.Annotations = make(map[string]string)
	qhObject.Annotations[LastAppliedConfigAnnotation] = qhObjectString

	quotahandler, err := c.handerLister.Get(app.Namespace, app.Name+"-"+"quotahandler")
	if err != nil {
		//log.Errorf("Get quotahandler for %s error : %s", app.Namespace+":"+app.Name+"-"+component.Name, err.Error())
		if errors.IsNotFound(err) {
			_, err = c.handlerClient.Create(qhObject)
			if err != nil {
				log.Errorf("Create quotahandler for %s error : %s", app.Namespace+":"+app.Name, err.Error())
			}
		}
	} else {

		if quotahandler != nil {
			if quotahandler.Annotations[LastAppliedConfigAnnotation] != qhObjectString {
				qhObject.ObjectMeta.ResourceVersion = quotahandler.ObjectMeta.ResourceVersion
				_, err = c.handlerClient.Update(qhObject)
				if err != nil {
					log.Errorf("Update quotahandler for %s error : %s", app.Namespace+":"+app.Name, err.Error())
				}
			}
		}
	}

	quotaruleObject := NewQuotaRuleObject(app)
	quotaruleObjectString := GetObjectApplied(quotaruleObject)
	quotaruleObject.Annotations = make(map[string]string)
	quotaruleObject.Annotations[LastAppliedConfigAnnotation] = quotaruleObjectString
	_, err = c.ruleLister.Get(app.Namespace, app.Name+"-"+"quotarule")
	if err != nil {
		//log.Errorf("Get quotarule for %s error : %s", app.Namespace+":"+app.Name+"-"+component.Name, err.Error())
		if errors.IsNotFound(err) {
			_, err = c.ruleClient.Create(&quotaruleObject)
			if err != nil {
				log.Errorf("Create quotarule for %s error : %s", app.Namespace+":"+app.Name, err.Error())
			}
		}
	}
	log.Infof("Sync quota config done for %s", app.Namespace)

	return nil
}

// sync trusted workload
func (c *controller) syncTrustedWorkload(component *v3.Component, app *v3.Application, ref *metav1.OwnerReference) error {
	resourceWorkloadType := "deployment"
	if resourceWorkloadType == "deployment" {
		deploy, err := c.deploymentLister.Get(app.Namespace, component.Name)
		if err != nil {
			log.Errorf("Get trusted deploy for %s error : %s", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
			return err
		}
		*ref = *(metav1.NewControllerRef(deploy, v1beta2.SchemeGroupVersion.WithKind("Deployment")))
		/*		ref.Name = deploy.Name
				ref.APIVersion = "apps/v1beta2"
				ref.Kind = "Deployment"
				ref.UID = deploy.ObjectMeta.UID*/
		object := deploy.DeepCopy()
		key := app.Name + "-" + component.Name + "-" + "workload"

		if val, _ := object.Spec.Template.Labels["app"]; val != key {
			object.Spec.Template.Labels["app"] = key
			newdeploy, err := c.deploymentClient.Update(object)
			if err != nil {
				log.Errorf("Update trusted deploy for %s error : %s", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
				return err
			}
			*ref = *(metav1.NewControllerRef(newdeploy, v1beta2.SchemeGroupVersion.WithKind("Deployment")))
		}
	}

	return nil
}

// zk update component state delete not exist version
func (c *controller) gc(namespace string, deletelist []string) (errlist []string) {
	for _, i := range deletelist {
		slices := strings.Split(i, "_")
		workloadname := slices[0] + "-" + slices[1] + "-" + "workload-" + slices[2]
		deletePolicy := metav1.DeletePropagationBackground
		err := c.deploymentClient.DeleteNamespaced(namespace, workloadname, &metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		if err != nil {
			log.Errorf("Delete Workload %s failed errinfo: %v", workloadname, err)
			errlist = append(errlist, i)
		}
	}
	return
}

// zk fucing
func (c *controller) syncFusing(podname, namespace string, set bool) {
	pod, err := c.podLister.Get(namespace, podname)
	if err != nil {
		log.Errorln("Get pod for namespace %s pod %s Error: %s", namespace, podname, err.Error())
	} else {
		if set {
			_, ok := pod.Labels["inpool"]
			if ok {
				log.Debugf("this pod %s already have this label", podname)
				return
			}
			pod.Labels["inpool"] = "yes"
			_, err = c.podClient.Update(pod)
			if err != nil {
				log.Errorf("Update pod %s for namespace %s Error: %s", podname)
			}
			return
		}
		_, ok := pod.Labels["inpool"]
		if ok {
			delete(pod.Labels, "inpool")
			_, err = c.podClient.Update(pod)
			if err != nil {
				log.Errorf("Update pod %s for namespace %s Error: %s", podname)
			}
		}
	}
}
