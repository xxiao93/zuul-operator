package zuul

import (
	"github.com/example-inc/zuul-operator/cmd/manager/tools/utils"
	cachev1alpha1 "github.com/example-inc/zuul-operator/pkg/apis/cache/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generatezuulexecutorVolumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      "zuul-executor-config",
			MountPath: "/etc/zuul/zuul.conf",
			SubPath: "zuul.conf",
			ReadOnly: utils.READONLY,
		},
		{
			Name:      "libzuul",
			MountPath: "/var/lib/zuul",
		},
	}
}

func generatezuulexecutorVolumes() []corev1.Volume {
	return []corev1.Volume{
		{
			Name: "zuul-executor-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "zuul-executor-config",
					},
				},
			},
		},
		{
			Name: "libzuul",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/lib/zuul",
				},
			},
		},
	}
}

//CreateZuulExecutorDeploymnet create zuulexecutor deployment
func CreateZuulExecutorDeployment(cr *cachev1alpha1.Zuul, serviceAccount *corev1.ServiceAccount) *appsv1.Deployment{
	labels := map[string]string{
		"k8s-app": "zuul-executor",
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "zuul-executor",
			Namespace: cr.ObjectMeta.Namespace,
			Labels:    labels,
		},

		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},

			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},

				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:            "zuul-executor",
							Image:           "hub.easystack.io/devops/ubuntu-source-zuul-executor:" + cr.Spec.ZuulVersion,
							ImagePullPolicy: "IfNotPresent",
							Command: []string{"sleep", "1d"},
							Env:          utils.GenerateEnvironmentVariables(cr),
							VolumeMounts: generatezuulexecutorVolumeMounts(),
							SecurityContext: &corev1.SecurityContext{RunAsUser: &zuul_user_id},
						},
					},
					Volumes:            generatezuulexecutorVolumes(),
					ServiceAccountName: serviceAccount.Name,
				},
			},
		},
	}
}

// CreateConfigMap - generate config map
func CreateZuulExecutorConfigMap(cr *cachev1alpha1.Zuul) *corev1.ConfigMap {

	templateInput := utils.TemplateInput{
		GerritServer:    cr.Spec.Gerrit.Server,
		GerritUser:      cr.Spec.Gerrit.User,
	}

	zuulconfig, _ := utils.GenerateConfig(templateInput, zuulexecutorconfigTemplate)

	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: "zuul-executor-config",
			Labels: map[string]string{
				"k8s-app": "zuul-executor",
			},
			Namespace: cr.ObjectMeta.Namespace,
		},

		Data: map[string]string{
			"zuul.conf": *zuulconfig,
		},
	}
}
