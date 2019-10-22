package zuul

import (
	"github.com/example-inc/zuul-operator/cmd/manager/tools/utils"
	cachev1alpha1 "github.com/example-inc/zuul-operator/pkg/apis/cache/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generatezuulmergerVolumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      "zuul-merger-config",
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

func generatezuulmergerVolumes() []corev1.Volume {
	return []corev1.Volume{
		{
			Name: "zuul-merger-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "zuul-merger-config",
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

//CreateZuulMergerDeploymnet create zuulmerger deployment
func CreateZuulMergerDeployment(cr *cachev1alpha1.Zuul, serviceAccount *corev1.ServiceAccount) *appsv1.Deployment{
	labels := map[string]string{
		"k8s-app": "zuul-merger",
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "zuul-merger",
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
							Name:            "zuul-merger",
							Image:           "hub.easystack.io/devops/ubuntu-source-zuul-merger:" + cr.Spec.ZuulVersion,
							ImagePullPolicy: "IfNotPresent",
							Command: []string{"sleep", "1d"},
							Env:          utils.GenerateEnvironmentVariables(cr),
							VolumeMounts: generatezuulmergerVolumeMounts(),
						},
					},
					Volumes:            generatezuulmergerVolumes(),
					ServiceAccountName: serviceAccount.Name,
				},
			},
		},
	}
}

// CreateConfigMap - generate config map
func CreateZuulMergerConfigMap(cr *cachev1alpha1.Zuul) *corev1.ConfigMap {

	templateInput := utils.TemplateInput{
		GerritServer:    cr.Spec.Gerrit.Server,
		GerritUser:      cr.Spec.Gerrit.User,
	}

	zuulconfig, _ := utils.GenerateConfig(templateInput, zuulmergerconfigTemplate)

	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: "zuul-merger-config",
			Labels: map[string]string{
				"k8s-app": "zuul-merger",
			},
			Namespace: cr.ObjectMeta.Namespace,
		},

		Data: map[string]string{
			"zuul.conf": *zuulconfig,
		},
	}
}
