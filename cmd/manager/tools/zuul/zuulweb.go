package zuul

import (
	"github.com/example-inc/zuul-operator/cmd/manager/tools/utils"
	cachev1alpha1 "github.com/example-inc/zuul-operator/pkg/apis/cache/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func generatezuulwebVolumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      "zuul-web-config",
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

func generatezuulwebVolumes() []corev1.Volume {
	return []corev1.Volume{
		{
			Name: "zuul-web-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "zuul-web-config",
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

//CreateZuulWebDeploymnet create zuulweb deployment
func CreateZuulWebDeployment(cr *cachev1alpha1.Zuul, serviceAccount *corev1.ServiceAccount) *appsv1.Deployment{
	labels := map[string]string{
		"k8s-app": "zuul-web",
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "zuul-web",
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
							Name:            "zuul-web",
							Image:           "hub.easystack.io/devops/ubuntu-source-zuul-web:" + cr.Spec.ZuulVersion,
							ImagePullPolicy: "IfNotPresent",
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									Name: "web-port",
									ContainerPort: 9001,
								},
							},
							Command: []string{"sleep", "1d"},
							Env:          utils.GenerateEnvironmentVariables(cr),
							VolumeMounts: generatezuulwebVolumeMounts(),
							SecurityContext: &corev1.SecurityContext{RunAsUser: &zuul_user_id},
						},
					},
					Volumes:            generatezuulwebVolumes(),
					ServiceAccountName: serviceAccount.Name,
				},
			},
		},
	}
}

// CreateConfigMap - generate config map
func CreateZuulWebConfigMap(cr *cachev1alpha1.Zuul) *corev1.ConfigMap {

	templateInput := utils.TemplateInput{}

	zuulconfig, _ := utils.GenerateConfig(templateInput, zuulwebconfigTemplate)

	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: "zuul-web-config",
			Labels: map[string]string{
				"k8s-app": "zuul-web",
			},
			Namespace: cr.ObjectMeta.Namespace,
		},

		Data: map[string]string{
			"zuul.conf": *zuulconfig,
		},
	}
}

// CreatezuulwebService generates zuul-web service
func CreateZuulWebService(cr *cachev1alpha1.Zuul) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      "zuul-web",
			Namespace: cr.ObjectMeta.Namespace,
		},

		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"k8s-app": "zuul-web",
			},
			Ports: []corev1.ServicePort{
				{
					Port: 9001,
					TargetPort: intstr.IntOrString{
						IntVal: int32(9001),
					},
				},
			},
		},
	}
}
