package utils

import (
	"fmt"
	"os"
	"text/template"

	nginxv1alpha1 "github.com/gokul-mylsami/nginx-operator/api/v1alpha1"
)

func NginxReload() {
	fmt.Printf("Reloading Nginx")
}

func NginxTemplateGenerator(Route nginxv1alpha1.NginxRoutes, templateFileName string) {
	ENV_TYPE := os.Getenv("ENV_TYPE")
	currentWorkingDir, _ := os.Getwd()
	var nginxTemplatePath string
	var templateFilePath string

	if ENV_TYPE == "PROD" {
		templateFilePath = "/etc/operator/templates/"
		nginxTemplatePath = "/etc/nginx/templates/"
	} else {
		nginxTemplatePath = currentWorkingDir + "/templates/"
		templateFilePath = currentWorkingDir + "/templates/"
	}

	// list all files in the template directory

	fileName := templateFilePath + templateFileName
	fmt.Println("Reading file: " + fileName)
	outputFileName := nginxTemplatePath + templateFileName

	tmpl, err := template.ParseFiles(fileName)
	if err != nil {
		fmt.Println("Error parsing template file : " + err.Error())
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Error creating output file : " + err.Error())
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, Route)
	if err != nil {
		fmt.Println("Error executing template file : " + err.Error())
	}

	fmt.Println("Template file generated" + outputFileName)
}
