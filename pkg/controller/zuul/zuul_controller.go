package zuul

import (
	"context"
	"strings"
	"time"

	"github.com/example-inc/zuul-operator/cmd/manager/tools"
	cachev1alpha1 "github.com/example-inc/zuul-operator/pkg/apis/cache/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_zuul")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Zuul Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileZuul{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("zuul-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Zuul
	err = c.Watch(&source.Kind{Type: &cachev1alpha1.Zuul{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Zuul
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cachev1alpha1.Zuul{},
	})
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cachev1alpha1.Zuul{},
	})
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cachev1alpha1.Zuul{},
	})

	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileZuul implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileZuul{}

// ReconcileZuul reconciles a Zuul object
type ReconcileZuul struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Zuul object and makes changes based on the state read
// and what is in the Zuul.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileZuul) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Zuul")

	// Fetch the Zuul instance
	instance := &cachev1alpha1.Zuul{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("Zuul CRD instance not found!")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	tool := tools.GetTools(instance)

	_, svcAccount, role, binding := tool.SetupAccountsAndBindings()

	existingSvcAccount := &corev1.ServiceAccount{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: svcAccount.Name, Namespace: svcAccount.Namespace}, existingSvcAccount)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Service Account")
		if err = createK8sObject(instance, svcAccount, r); err != nil {
			return reconcile.Result{}, err
		}
		return requeAfter(1, nil)
	}

	existingRole := &rbacv1.ClusterRole{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: role.Name}, existingRole)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Cluster Role")
		if err = createK8sObject(instance, role, r); err != nil {
			return reconcile.Result{}, err
		}
		return requeAfter(1, nil)
	}

	existingBinding := &rbacv1.ClusterRoleBinding{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: binding.Name}, existingBinding)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Role Binding")
		if err = createK8sObject(instance, binding, r); err != nil {
			return reconcile.Result{}, err
		}
		return requeAfter(1, nil)
	}

	existingZuulSchedulerConfigMap, configMap := tool.ZuulScheduler.GetZuulSchedulerConfigMap()
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, existingZuulSchedulerConfigMap)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating ZuulScheduler Config Map")
		if err = createK8sObject(instance, configMap, r); err != nil {
			return reconcile.Result{}, err
		}
		return requeAfter(5, nil)
	}

	existingZuulScheduler, zuulscheduler := tool.ZuulScheduler.GetZuulSchedulerDeployment()
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: zuulscheduler.Name, Namespace: zuulscheduler.Namespace}, existingZuulScheduler)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating ZuulScheduler")
		if err = createK8sObject(instance, zuulscheduler, r); err != nil {
			return reconcile.Result{}, err
		}
		return requeAfter(5, nil)
	}

	/*
		// Updation
	*/

	// Gerrit info Update, it should put first because it has delete action !!!
	gerrit_server := instance.Spec.Gerrit.Server
	gerrit_port := instance.Spec.Gerrit.Port

	env := existingZuulScheduler.Spec.Template.Spec.Containers[0].Env
	env_map := generateEnvMap(env)

	// If update gerrit message, we need recreate configmap and deployment
	if ( env_map["gerrit_server"] != gerrit_server || env_map["gerrit_port"] != gerrit_port ) {
		reqLogger.Info("Begin Delete ZuulScheduler Configmap")
		if err = r.client.Delete(context.TODO(), existingZuulSchedulerConfigMap); err != nil {
			return reconcile.Result{}, err
		}
		reqLogger.Info("Begin Delete ZuulScheduler Deployment")
		if err = r.client.Delete(context.TODO(), existingZuulScheduler); err != nil {
			return reconcile.Result{}, err
		}
		return requeAfter(5, nil)
	}

	// ZuulScheduler Deployment Size Update
	size := instance.Spec.ZuulScheduler.Size
	if *existingZuulScheduler.Spec.Replicas != size {
		existingZuulScheduler.Spec.Replicas = &size
		if err = r.client.Update(context.TODO(), existingZuulScheduler); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}

	// Zuul Version update
	version := instance.Spec.ZuulVersion
	actual_zc_image_prefix, actual_zc_image_version := generateImageDetail(existingZuulScheduler)

	if actual_zc_image_version != version {
		expect_image := strings.Join([]string{actual_zc_image_prefix, version}, ":")
		existingZuulScheduler.Spec.Template.Spec.Containers[0].Image = expect_image

		if err = r.client.Update(context.TODO(), existingZuulScheduler); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}

	return reconcile.Result{}, nil
}

func generateEnvMap(env []corev1.EnvVar) (map[string]string) {
	dict := make(map[string]string)

	for _, v := range env {
		dict[strings.ToLower(v.Name)] = v.Value
	}
	return dict
}

func generateImageDetail(obj metav1.Object) (string, string) {
	var image string

	switch t := obj.(type) {
	case *appsv1.Deployment:
		image = t.Spec.Template.Spec.Containers[0].Image
	}
	return strings.Split(image, ":")[0], strings.Split(image, ":")[1]
}

func createK8sObject(instance *cachev1alpha1.Zuul, obj metav1.Object, r *ReconcileZuul) error {
	var err error
	err = controllerutil.SetControllerReference(instance, obj, r.scheme)

	if err != nil {
		return err
	}

	switch t := obj.(type) {
	case *corev1.ServiceAccount:
		err = r.client.Create(context.TODO(), t)
	case *rbacv1.ClusterRole:
		err = r.client.Create(context.TODO(), t)
	case *rbacv1.ClusterRoleBinding:
		err = r.client.Create(context.TODO(), t)
	case *corev1.Namespace:
		err = r.client.Create(context.TODO(), t)
	case *corev1.ConfigMap:
		err = r.client.Create(context.TODO(), t)
	case *corev1.Service:
		err = r.client.Create(context.TODO(), t)
	case *appsv1.DaemonSet:
		err = r.client.Create(context.TODO(), t)
	case *appsv1.Deployment:
		err = r.client.Create(context.TODO(), t)
	}
	return err
}

func requeAfter(sec int, err error) (reconcile.Result, error) {
	t := time.Duration(sec)
	return reconcile.Result{RequeueAfter: time.Second * t}, err
}

