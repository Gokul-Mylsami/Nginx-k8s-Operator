# 🚀 NGINX Operator  

## 📌 Overview  
The **NGINX Operator** is a custom **CRD tool** that can be used in **Kubernetes** as an alternative to vanilla **NGINX**.  

### 🔍 Key Differences from Vanilla NGINX  
✅ With **vanilla NGINX**, configuration files must be **manually mounted** or managed through other methods. If any changes are made, a **rollout restart** of NGINX components is required for the updates to take effect.  

✅ The **NGINX Operator** eliminates this issue by using custom CRDs such as **`NginxRoutes`** and **`NginxUpstreams`** to **dynamically mount configurations**. This ensures that NGINX picks up new changes **without requiring a reload**.  

### 🎯 Benefits of Using NGINX Operator  
🔹 **No manual configuration mounting**  
🔹 **No need for rollout restarts**  
🔹 **Dynamic and seamless config updates**  

🔥 **Simplify your NGINX management with the NGINX Operator!** 🚀  


## 🚀 Getting Started  

### Step 1: Install the NGINX Operator CRDs  

Run the following commands to install the required CRDs for the NGINX Operator:  

```sh
kubectl apply -f https://raw.githubusercontent.com/Gokul-Mylsami/Nginx-k8s-Operator/refs/heads/main/config/crd/bases/nginx.gokul-mylsami.com_nginxroutes.yaml

kubectl apply -f https://raw.githubusercontent.com/Gokul-Mylsami/Nginx-k8s-Operator/refs/heads/main/config/crd/bases/nginx.gokul-mylsami.com_nginxupstreams.yaml

```

### Step 2: Create a Route Template
Define a route template based on your use case.


```sh 
# Example Template 
server {
    include /etc/nginx/conf.d/*.conf;
    
    listen {{ .Spec.ServerPort }};

    {{- if .Spec.ServerName }}
    server_name {{ .Spec.ServerName }};
    {{- end }}

    {{- if .Spec.CustomDirectives }}
        {{- range .Spec.CustomDirectives }}
            {{ . }}
        {{- end }}
    {{- end }}

    {{- if .Spec.TLSCertificate }}
        ssl_certificate /etc/nginx/ssl/{{.Spec.TLSCertificate.Name}}-{{.Spec.TLSCertificate.Namespace}}.crt;
        ssl_certificate_key /etc/nginx/ssl/{{.Spec.TLSCertificate.Name}}-{{.Spec.TLSCertificate.Namespace}}.key;
    {{- end }}

    {{- if .Spec.CustomLocations }}
        {{- range .Spec.CustomLocations }}
            location {{ .Location }} {
                {{ .Definition }}
            }
        {{- end }}
    {{- end }}
}

```

This NGINX template must be provided when creating routes to generate an NGINX configuration tailored to the application's needs.

**Note**: The NGINX Operator uses Go templating functions, so ensure the syntax is correct. For more details, refer to the <a href="https://pkg.go.dev/text/template">Go templating documentation</a>.

### Step 3: Deploy the NGINX Operator

Create a Deployment or StatefulSet for the NGINX Operator using the image: `gokulmylsami/operator:v3`

Example: 

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-operator
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      serviceAccountName: nginx-operator-sa
      containers:
        - name: nginx
          image: gokulmylsami/operator:v3
          imagePullPolicy: Always
          resources:
            limits:
              memory: "512Mi"
              cpu: "500m"
            requests:
              memory: "256Mi"
              cpu: "250m"
          env:
            - name: ENV_TYPE
              value: PROD
          volumeMounts:
            - name: nginx-config-volume
              mountPath: /etc/operator/templates/
          ports:
            - containerPort: 80
      volumes:
        - name: nginx-config-volume
          configMap:
            name: nginx-config   # This is the template ConfigMap created in Step 2  

```

### Step 4: Create and Deploy a Route

Define and deploy a NginxRoute resource in your Kubernetes cluster.

```yaml
apiVersion: nginx.gokul-mylsami.com/v1alpha1
kind: NginxRoutes
metadata:
  labels:
    app.kubernetes.io/name: operator
    app.kubernetes.io/managed-by: kustomize
  name: nginxroutes-sample
spec:
  serverPort: 80
  serverName: "_"
  templateFile: "main.conf.template"
  customLocations: 
    - location: "/"
      definition: | 
        return 200 'Hello, World!';
```

🎉 You have successfully set up the NGINX Operator! 🚀
Now, your NGINX configurations can be dynamically updated without requiring a full reload.

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

