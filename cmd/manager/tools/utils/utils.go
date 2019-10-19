package utils

import (
	"bytes"
	"text/template"

	cachev1alpha1 "github.com/example-inc/zuul-operator/pkg/apis/cache/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

var READONLY = true

// Setup Zuul Env
func GenerateEnvironmentVariables(cr *cachev1alpha1.Zuul) []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name:  "GERRIT_SERVER",
			Value: cr.Spec.Gerrit.Server,
		},
		{
			Name:  "GERRIT_USER",
			Value: cr.Spec.Gerrit.User,
		},
		{
			Name:  "GERRIT_PORT",
			Value: cr.Spec.Gerrit.Port,
		},
	}
}

// TemplateInput defines the input template placeholder
type TemplateInput struct {
	GerritServer  string
	GerritUser    string
}

func GenerateConfig(input TemplateInput, templateFile string) (*string, error) {
	output := new(bytes.Buffer)
	tmpl, err := template.New("config").Parse(templateFile)
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(output, input)
	if err != nil {
		return nil, err
	}
	outputString := output.String()
	return &outputString, nil
}
