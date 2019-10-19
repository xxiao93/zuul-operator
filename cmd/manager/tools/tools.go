package tools

import (
	"github.com/example-inc/zuul-operator/cmd/manager/tools/zuul"
	cachev1alpha1 "github.com/example-inc/zuul-operator/pkg/apis/cache/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Tools structure declarations
type Tools struct {
	cr            *cachev1alpha1.Zuul
	ZuulScheduler ZuulScheduler
}

func (t *Tools) init() {
	t.ZuulScheduler = ZuulScheduler{
		cr: t.cr,
	}
}

// SetupAccountsAndBindings creates Service Account, Cluster Role and Cluster Binding for Zuul
func (t *Tools) SetupAccountsAndBindings() (*corev1.Namespace, *corev1.ServiceAccount, *rbacv1.ClusterRole, *rbacv1.ClusterRoleBinding) {
	namespace := t.createNamespace()
	svcAccount := t.createServiceAccount()
	clusterRole := t.createClusterRole()
	roleBinding := t.createRoleBinding(clusterRole, svcAccount)

	t.ZuulScheduler.serviceAccount = svcAccount
	return namespace, svcAccount, clusterRole, roleBinding
}

// ZuulScheduler structure
type ZuulScheduler struct {
	cr             *cachev1alpha1.Zuul
	serviceAccount *corev1.ServiceAccount
}

// GetConfigMap returns ZuulScheduler ConfigMap
func (z *ZuulScheduler) GetZuulSchedulerConfigMap() (*corev1.ConfigMap, *corev1.ConfigMap) {
	return &corev1.ConfigMap{}, zuul.CreateZuulSchedulerConfigMap(z.cr)
}

// GetDeployment returns ZuulScheduler Deployment
func (z *ZuulScheduler) GetZuulSchedulerDeployment() (*appsv1.Deployment, *appsv1.Deployment) {
	return &appsv1.Deployment{}, zuul.CreateZuulSchedulerDeployment(z.cr, z.serviceAccount)
}

/* -------------------------------
// Util Functions
// ------------------------------- */
func (t Tools) createServiceAccount() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      "zuul",
			Namespace: t.cr.ObjectMeta.Namespace,
		},
	}
}

func (t Tools) createClusterRole() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1beta1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: "zuul",
		},

		Rules: []rbacv1.PolicyRule{{
			APIGroups: []string{"*"},
			Resources: []string{"*"},
			Verbs:     []string{"*"},
		}},
	}
}

func (t Tools) createRoleBinding(clusterRole *rbacv1.ClusterRole, svcAccount *corev1.ServiceAccount) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1beta1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: "zuul",
		},

		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     clusterRole.TypeMeta.Kind,
			Name:     clusterRole.ObjectMeta.Name,
		},

		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      svcAccount.ObjectMeta.Name,
			Namespace: t.cr.ObjectMeta.Namespace,
		}},
	}
}

func (t Tools) createNamespace() *corev1.Namespace {
	return &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: t.cr.ObjectMeta.Namespace,
		},
	}
}

// ------------------------------

// GetTools returns an instance of Tools
func GetTools(customResource *cachev1alpha1.Zuul) *Tools {
	tools := Tools{
		cr: customResource,
	}
	tools.init()
	return &tools
}
