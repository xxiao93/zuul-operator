package zuul

import (
	"github.com/example-inc/zuul-operator/cmd/manager/tools/utils"
	cachev1alpha1 "github.com/example-inc/zuul-operator/pkg/apis/cache/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ro = true

func generateVolumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      "zuul-scheduler-config",
			MountPath: "/etc/zuul/zuul.conf",
			SubPath: "zuul.conf",
			ReadOnly: ro,
		},
		{
			Name:      "zuul-scheduler-config",
			MountPath: "/var/lib/zuul/tenant-config/main.yaml",
			SubPath: "main.yaml",
			ReadOnly: ro,
		},
		{
			Name:      "libzuul",
			MountPath: "/var/lib/zuul",
		},
	}
}

func generateVolumes() []corev1.Volume {
	return []corev1.Volume{
		{
			Name: "zuul-scheduler-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "zuul-scheduler-config",
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

//CreateZuulSchedulerDeploymnet create zuulscheduler deployment
func CreateZuulSchedulerDeployment(cr *cachev1alpha1.Zuul, serviceAccount *corev1.ServiceAccount) *appsv1.Deployment{
	labels := map[string]string{
		"k8s-app": "zuul-scheduler",
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "zuul-scheduler",
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
							Name:            "zuul-scheduler",
							Image:           "hub.easystack.io/devops/ubuntu-source-zuul-scheduler:" + cr.Spec.ZuulVersion,
							ImagePullPolicy: "IfNotPresent",
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									Name: "scheduler-port",
									ContainerPort: 2020,
								},
							},
							Command: []string{"sleep", "1d"},
							Env:          utils.GenerateEnvironmentVariables(cr),
							VolumeMounts: generateVolumeMounts(),
						},
					},
					Volumes:            generateVolumes(),
					ServiceAccountName: serviceAccount.Name,
				},
			},
		},
	}
}

// CreateConfigMap - generate config map
func CreateZuulSchedulerConfigMap(cr *cachev1alpha1.Zuul) *corev1.ConfigMap {

	templateInput := utils.TemplateInput{
		GerritServer:    cr.Spec.Gerrit.Server,
		GerritUser:      cr.Spec.Gerrit.User,
	}

	zuulconfig, _ := utils.GenerateConfig(templateInput, zuulschedulerconfigTemplate)
	mainyamlconfig, _ := utils.GenerateConfig(templateInput, mainyamlTemplate)

	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: "zuul-scheduler-config",
			Labels: map[string]string{
				"k8s-app": "zuul-scheduler",
			},
			Namespace: cr.ObjectMeta.Namespace,
		},

		Data: map[string]string{
			"main.yaml": *mainyamlconfig,
			"zuul.conf": *zuulconfig,
		},
	}
}
