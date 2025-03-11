package utils

import (
	"fmt"
	"os"
	"text/template"

	nginxv1alpha1 "github.com/gokul-mylsami/nginx-operator/api/v1alpha1"
)

func UpstreamTemplateGenerator(Route nginxv1alpha1.NginxUpstream, templateFileName string) {
	ENV_TYPE := os.Getenv("ENV_TYPE")
	currentWorkingDir, _ := os.Getwd()
	var nginxTemplatePath string
	var templateFilePath string

	if ENV_TYPE == "PROD" {
		templateFilePath = "/etc/operator/templates/"
		nginxTemplatePath = "/etc/nginx/conf.d/"
	} else {
		nginxTemplatePath = currentWorkingDir + "/"
		templateFilePath = currentWorkingDir + "/templates/"
	}
	// Generate the nginx configuration file
	outputTemplateFileName := nginxTemplatePath + Route.Name + ".conf"
	templateFileName = templateFilePath + templateFileName

	tmpl, err := template.ParseFiles(templateFileName)
	if err != nil {
		fmt.Println("[error] Error parsing the template file : " + err.Error())
	}

	// write the template to a file and remove .template extension
	file, err := os.Create(outputTemplateFileName)
	if err != nil {
		fmt.Println("[error] Error creating the template file : " + err.Error())
	}

	err = tmpl.Execute(file, Route)
	if err != nil {
		fmt.Println("[error] Error executing the template : " + err.Error())
	}

	fmt.Println("[info] Upstream Template file generated : " + outputTemplateFileName)
}
