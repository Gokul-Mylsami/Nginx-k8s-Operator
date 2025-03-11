package utils

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	nginxv1alpha1 "github.com/gokul-mylsami/nginx-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NginxReload() error {
	// First check if nginx is running
	checkCmd := exec.Command("/bin/sh", "-c", "nginx -t")
	if output, err := checkCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("nginx configuration test failed: %v\nOutput: %s", err, output)
	}

	// Reload nginx
	reloadCmd := exec.Command("/bin/sh", "-c", "nginx -s reload")
	if output, err := reloadCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("nginx reload failed: %v\nOutput: %s", err, output)
	}

	fmt.Println("[info] Nginx configuration reloaded successfully")
	return nil
}

func SecretGenerator(r client.Client, ctx context.Context, Route nginxv1alpha1.NginxRoutes) {
	secretName := Route.Spec.TLSCertificate.Name
	secretNamespace := Route.Spec.TLSCertificate.Namespace
	if secretName != "" || secretNamespace != "" {
		secret := &corev1.Secret{}
		err := r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: secretNamespace}, secret)
		if err != nil {
			fmt.Println("[error] Error fetching secret : " + err.Error())
		}

		// Write the secret to a file
		secretData := secret.Data
		// decode the secret data
		tlsCertificate := string(secretData["tls.crt"])
		tlsKey := string(secretData["tls.key"])

		// create ssl directory if it does not exist
		if _, err := os.Stat("/etc/nginx/ssl"); os.IsNotExist(err) {
			os.Mkdir("/etc/nginx/ssl", 0755)
		}

		// write the tls certificate and key to a file
		tlsCertificateFile, err := os.Create("/etc/nginx/ssl/" + secretName + "-" + secretNamespace + ".crt")
		if err != nil {
			fmt.Println("[error] Error creating tls certificate file : " + err.Error())
		}

		tlsKeyFile, err := os.Create("/etc/nginx/ssl/" + secretName + "-" + secretNamespace + ".key")
		if err != nil {
			fmt.Println("[error] Error creating tls key file : " + err.Error())
		}

		tlsCertificateFile.WriteString(tlsCertificate)
		tlsKeyFile.WriteString(tlsKey)
	} else {
		fmt.Println("[info] No TLS Certificate found for the route : " + Route.Name)
	}
}

func NginxTemplateGenerator(Route nginxv1alpha1.NginxRoutes, templateFileName string) {
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

	// list all files in the template directory

	fileName := templateFilePath + templateFileName
	fmt.Println("Reading file: " + fileName)
	// remove .template from templateFileName
	outputFileName := nginxTemplatePath + templateFileName[:len(templateFileName)-9]

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

	// Reload Nginx
	NginxReload()
}
